package models

import org.ktorm.entity.Entity
import org.ktorm.schema.*
import java.time.LocalDateTime
import java.util.UUID

interface User: Entity<User> {
    val id: UUID
    var passwordHash: String
    var verified: Boolean
    val modifiedAt: LocalDateTime
    val createdAt: LocalDateTime
}

object Users: Table<User>("m_user") {
    val id = uuid("id").primaryKey().bindTo { it.id }
    val passwordHash = varchar("password_hash").bindTo { it.passwordHash }
    val verified = boolean("verified").bindTo { it.verified }
    val createdAt = datetime("created_at").bindTo { it.createdAt }
    val modifiedAt = datetime("modified_at").bindTo { it.modifiedAt }
}