package models

import java.time.LocalDateTime
import java.util.UUID

class Account(
    val id: UUID,
    val password_hash: String,
    val verified: Boolean,
    val mfa: Any,
    val created_at: LocalDateTime,
    val updated_at: LocalDateTime
) {

}