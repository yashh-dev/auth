import at.favre.lib.crypto.bcrypt.BCrypt
import at.favre.lib.crypto.bcrypt.LongPasswordStrategy


class Crypto {
    companion object {
        fun hash(password: String): String {
            return BCrypt.with(LongPasswordStrategy.PassThroughStrategy(),).hashToString(13, password.toCharArray())
        }


        fun verify(password: String, hash: String): Boolean {
            return BCrypt.verifyer().verifyStrict(password.toCharArray(), hash.toCharArray()).verified
        }
    }

}