import com.google.gson.Gson
import com.rabbitmq.client.AMQP
import com.rabbitmq.client.Delivery;
import org.postgresql.ds.PGPoolingDataSource
import javax.sql.DataSource
import org.jdbi.v3.core.Jdbi
import org.json.JSONObject
import kotlin.jvm.optionals.getOrNull

class Service(private val name: String, connectUri: String) {
    private val rabbitService: RabbitService = RabbitService()
    private val map: HashMap<String, (Delivery) -> Unit> = HashMap()
    private val database: Jdbi = Jdbi.create(setupDataSource(connectUri))

    fun add(event: String, func: (Delivery) -> Unit) {
        map[event] = func
    }

    private fun setupDataSource(connectUri: String): DataSource {
        val source: PGPoolingDataSource = PGPoolingDataSource()
        source.initialConnections = 3
        source.maxConnections = 10
        source.serverNames = arrayOf("192.168.1.28")
        source.databaseName = "miauw"
        source.user = "miauw_user"
        source.password = "miauw_password"
        return source
    }

    fun start() {
        for ((event, handler) in map) {
            rabbitService.declareQueue("${name}.${event}")
            rabbitService.start("${name}.${event}", handler)
        }
    }

    fun handleUserCreate(delivery: Delivery) {
        val data = Gson().fromJson(String(delivery.body), UserCreateDTO::class.java)
        val passwordHash: String = Crypto.hash(data.password as String)
        database.inTransaction<Any, Exception> { handler ->
            handler.execute("insert into accounts(\"id\", \"password_hash\") values('${data.id}','$passwordHash');")
        }
        val verificationToken = JWToken.createVerificationToken(data.id as String)
        rabbitService.sendVerificationEMail(data.id as String, verificationToken)
    }

    fun handleUserLogin(delivery: Delivery): String {
        val data: UserLoginDTO = Gson().fromJson(String(delivery.body), UserLoginDTO::class.java)
        var response = ""
        val account: Map<String, Any?>? = database.withHandle<Map<String, Any?>, Exception> {
            it.createQuery("select * from accounts where id::text = '${data.id}'").mapToMap().findOne().getOrNull()
        }
        if (account == null) {
            response = JSONObject(
                mapOf<String, Any>(
                    "type" to "https://auth.miauw.social/login/account-not-found",
                    "title" to "The account is not found!",
                    "detail" to "Your account link with the profile is not found. This should not happen. Please reach out to admin.",
                    "status" to 404
                )
            ).toString()
        } else if (account["verified"] != true) {
            response = JSONObject(
                mapOf<String, Any>(
                    "type" to "https://auth.miauw.social/login/account-not-verified",
                    "title" to "The account is not verified!",
                    "detail" to "Your account is is not verified, therefore login failed.",
                    "status" to 403
                )
            ).toString()
        }
        if (!Crypto.verify(data.password as String, account?.get("password_hash") as String)) {
            response = JSONObject(
                mapOf<String, Any>(
                    "type" to "https://auth.miauw.social/login/wrong-password",
                    "title" to "The password is wrong!",
                    "detail" to "Your provided password does not match the password hash in the database.",
                    "status" to 403
                )
            ).toString()
        } else {
            val userId: String = account["id"].toString()
            val x = database.inTransaction<Map<String, Any?>, Exception> {
                it.createQuery("insert into sessions(\"account\") values('${userId}') returning id;")
                    .mapToMap()
                    .one()
            }
            response = x["id"].toString()
        }

        return response;
    }
}


fun main(args: Array<String>) {
    val service = Service("auth", "jdbc:postgres://miauw_user:miauw_password@192.168.1.28/miauw")
    service.add("password.initial", service::handleUserCreate)
    service.add("login", service::handleUserLogin)
    service.start()
}