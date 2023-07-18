import com.rabbitmq.client.CancelCallback
import com.rabbitmq.client.ConnectionFactory
import com.rabbitmq.client.DeliverCallback
import com.rabbitmq.client.Delivery

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
        channel.queueDeclare(queueName, true, false, true, emptyMap())
        channel.close()
        newConnection.close()
        return this
    }

    fun start(event: String, handler: (String, Delivery) -> Unit) {
        val connection = getFactory().newConnection()
        val channel = connection.createChannel()

        val cancelCallback = CancelCallback { consumerTag: String? -> println("Cancelled... $consumerTag") }

        channel.basicConsume(
            event,
            false,
            handler,
            cancelCallback
        )
    }
}