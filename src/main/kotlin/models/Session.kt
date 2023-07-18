package models

import org.ktorm.entity.Entity
import org.ktorm.schema.Table
import org.ktorm.schema.datetime
import org.ktorm.schema.uuid
import java.time.LocalDateTime
import java.util.UUID


interface Session: Entity<Session> {
    companion object: Entity.Factory<Session>()
    val id: UUID
    var user: User
    var createdAt: LocalDateTime
}
object Sessions: Table<Session>("m_session") {
    val id = uuid("id").primaryKey().bindTo { it.id }
    val user = uuid("user_id").references(Users){ it.user }
    val createdAt = datetime("created_at").bindTo { it.createdAt }
}