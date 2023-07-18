import org.json.JSONObject
import org.jose4j.jwa.AlgorithmConstraints
import org.jose4j.jws.AlgorithmIdentifiers
import org.jose4j.jws.JsonWebSignature
import org.jose4j.jwt.JwtClaims
import org.jose4j.jwt.consumer.ErrorCodes
import org.jose4j.jwt.consumer.InvalidJwtException
import org.jose4j.jwt.consumer.JwtConsumer
import org.jose4j.jwt.consumer.JwtConsumerBuilder
import org.jose4j.keys.AesKey

class JWToken {

    companion object {
        fun createVerificationToken(userId: String): String {
            val aesKey = AesKey(System.getenv("JWT_SECRET").toByteArray())
            // claims
            val claims = JwtClaims()
            claims.issuer = "https://auth.miauw.social"
            claims.audience = mutableListOf(userId)
            claims.setExpirationTimeMinutesInTheFuture(180F)
            claims.setGeneratedJwtId()
            claims.setIssuedAtToNow()
            claims.setNotBeforeMinutesInThePast(2F)
            claims.subject = userId
            // jws
            val jws: JsonWebSignature = JsonWebSignature()
            jws.payload = claims.toJson()
            jws.key = aesKey
            jws.algorithmHeaderValue = AlgorithmIdentifiers.HMAC_SHA256
            // jwt
            return jws.compactSerialization
        }

        fun verifyToken(token: String): Any {
            val aesKey = AesKey(System.getenv("JWT_SECRET").toByteArray())
            val consumer: JwtConsumer = JwtConsumerBuilder()
                .setRequireJwtId()
                .setRequireIssuedAt()
                .setRequireSubject()
                .setAllowedClockSkewInSeconds(60)
                .setExpectedIssuer("https://auth.miauw.social")
                .setVerificationKey(aesKey)
                .setJwsAlgorithmConstraints(
                    AlgorithmConstraints.ConstraintType.PERMIT,
                    AlgorithmIdentifiers.HMAC_SHA256
                )
                .build()
            try {
                return consumer.processToClaims(token)
            } catch (e: InvalidJwtException) {
                if (e.hasExpired()) {
                    return JSONObject(
                        mapOf(
                            "type" to "https://auth.miauw.social/verify/token-expired",
                            "title" to "Your token expired!",
                            "detail" to "The provided token has expired and is therefore not valid.",
                            "status" to 401
                        )
                    )
                } else if (e.hasErrorCode(ErrorCodes.ISSUER_INVALID)) {
                    return JSONObject(
                        mapOf(
                            "type" to "https://auth.miauw.social/verify/token-issuer-invalid",
                            "title" to "Your token is invalid!",
                            "detail" to "The issuer of your token could not be valided.",
                            "status" to 401
                        )
                    )
                } else if (e.hasErrorCode(ErrorCodes.JWT_ID_MISSING)) {
                    return JSONObject(
                        mapOf(
                            "type" to "https://auth.miauw.social/verify/token-jid-invalid",
                            "title" to "Your token is invalid!",
                            "detail" to "The provided token is missing its identifier.",
                            "status" to 401
                        )
                    )
                } else {
                     return JSONObject(
                        mapOf(
                            "type" to "https://auth.miauw.social/verify/token-jid-invalid",
                            "title" to "Your token is invalid!",
                            "detail" to "The server was not able to verify your token. (Reason: ${e.errorDetails}",
                            "status" to 401
                        )
                    )
                }
            }
        }
    }


}