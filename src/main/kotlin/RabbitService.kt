import com.rabbitmq.client.*
import org.json.JSONObject


class RabbitService {
    private val connectionFactory: ConnectionFactory = ConnectionFactory()

    init {
        connectionFactory.host = "192.168.1.28"
        connectionFactory.port = 5672
        connectionFactory.username = "guest"
        connectionFactory.password = "guest"
    }

    fun getFactory(): ConnectionFactory {
        return connectionFactory
    }

    fun declareQueue(queueName: String): RabbitService {
        val newConnection = getFactory().newConnection()
        val channel = newConnection.createChannel()
        channel.queueDeclare(queueName, true, false, false, emptyMap())
        channel.close()
        newConnection.close()
        return this
    }

    fun start(event: String, handler: (Delivery) -> Unit, remoteProcedureCall: Boolean = false) {
        val connection = getFactory().newConnection()
        val channel = connection.createChannel()

        val cancelCallback = CancelCallback { consumerTag: String? -> println("Cancelled... $consumerTag") }
        val deliverCallback = DeliverCallback { consumerTag: String?, delivery: Delivery ->
            val replyProps = AMQP.BasicProperties.Builder()
                .correlationId(delivery.properties.correlationId)
                .build()
            var response: String = ""
            try {
                response = handler(delivery).toString()
            } catch (e: RuntimeException) {
                println(" [.] $e")
            } finally {
                if (remoteProcedureCall) {
                    channel.basicPublish(
                        "",
                        delivery.properties.replyTo,
                        replyProps,
                        response.toByteArray(charset("UTF-8"))
                    )
                }
                channel.basicAck(delivery.envelope.deliveryTag, false)
            }
        }
        channel.basicConsume(
            event,
            false,
            deliverCallback,
            cancelCallback
        )
    }

    fun sendVerificationEMail(recipient: String, vid: String) {
        val conn = getFactory().newConnection()
        val channel = conn.createChannel()
        val content = JSONObject(
            mapOf(
                "type" to "sign_up",
                "recipient" to recipient,
                "payload" to JSONObject(
                    mapOf(
                        "vid" to vid
                    )
                )
            )
        )

        channel.basicPublish("", "email", null, content.toString().toByteArray())
    }
}