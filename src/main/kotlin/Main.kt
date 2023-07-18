import com.google.gson.Gson
import com.rabbitmq.client.AMQP
import com.rabbitmq.client.Delivery;
import models.User
import org.ktorm.database.Database
import org.postgresql.ds.PGPoolingDataSource
import javax.sql.DataSource

class Service(private val name: String, connectUri: String) {
    private val rabbitService: RabbitService = RabbitService()
    private val map: HashMap<String, (String, Delivery) -> Unit> = HashMap()
    val database: Database = Database.connect(setupDataSource(connectUri))


    fun add(event: String, func: (String, Delivery) -> Unit) {
        map[event] = func
    }

    private fun setupDataSource(connectUri: String): DataSource {
        val source: PGPoolingDataSource = PGPoolingDataSource()
        source.initialConnections = 3
        source.maxConnections = 10
        source.serverNames = arrayOf("192.168.1.28")
        source.databaseName = "miauw"
        source.user = "miauw_user"
        source.password ="miauw_password"
        return source
    }

    fun start() {
        for ((event, handler) in map) {
            rabbitService.declareQueue("${name}.${event}")
            rabbitService.start("${name}.${event}", handler)
        }
    }
}


fun handleUserCreate(consumerTag: String, delivery: Delivery): String {
    val data = Gson().fromJson(String(delivery.body), UserCreateDTO::class.java)
    val passwordHash: String = Crypto.hash(data.password as String)

    return JWToken.createVerificationToken(data.id as String)
}

fun handleUserLogin(consumerTag: String, delivery: Delivery) {
    val user: UserLoginDTO = Gson().fromJson(String(delivery.body), UserLoginDTO::class.java)
    val replyProps: AMQP.BasicProperties =
        AMQP.BasicProperties.Builder().correlationId(delivery.properties.correlationId).build()
    try {

    } catch (e: RuntimeException) {
        println(e)
    } finally {
    }
}

fun main(args: Array<String>) {
    val service = Service("auth", "jdbc:postgres://miauw_user:miauw_password@192.168.1.28/miauw")
    service.add("test1", ::handleUserCreate)
    service.add("test2", ::handleUserLogin)
    service.start()
}