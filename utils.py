from redis.asyncio import BlockingConnectionPool, Redis
import uuid
from jwt import JWT, jwk_from_pem
from jwt.utils import get_int_from_datetime
import uuid
import datetime

class Rediser:
    def __init__(self, url: str, db: int = 0) -> None:
        self.url = url
        self.__pool = BlockingConnectionPool(host="192.168.1.28")
        self.radis: Redis = None
        self.db: int = db

    async def init(self) -> "Rediser":
        self.radis = Redis(
            host="192.168.1.28",
            connection_pool=self.__pool,
            auto_close_connection_pool=True,
            db=self.db,
        )
        return self

    async def set(self, key: str, value: str) -> None:
        await self.radis.set(key, value)

    async def get(self, key: str):
        """gets value from key"""
        return await self.radis.get(key)

    async def set_w_ttl(self, key: str, value: str, ttl: int):
        await self.radis.set(key, value, ex=ttl)


class SessionRediser(Rediser):
    def __init__(self, url: str):
        super().__init__(url, 0)

    async def exists(self, sid: str) -> bool:
        """checks if session the user is trying to use exists"""
        uid = await self.get(sid)
        if not uid:
            # session doesnt exist
            return False
        return True

    async def new(self, sid: uuid.UUID, uid: uuid.UUID, ttl: int):
        await self.radis.set(str(sid), str(uid), ex=172800)



instance = JWT()


def create_verification_token(user_id: uuid.UUID) -> str:
    """creates verification token from userid"""
    payload = {
        "iss": "miauw.social/auth",
        "sub": str(user_id),
        "iat": get_int_from_datetime(datetime.datetime.now()),
        "exp": get_int_from_datetime(
            datetime.datetime.now() + datetime.timedelta(minutes=15)
        ),
    }
    with open("certs/private.pem", "rb") as pk:
        signing_key: str = jwk_from_pem(pk.read())
    return instance.encode(payload, signing_key, alg="RS256")


def verify_verification_token(token: str) -> dict[str, str]:
    """verifies token"""
    with open("certs/public.pem", "rb") as pk:
        verifying_key = jwk_from_pem(pk.read())
    return instance.decode(token, verify_verification_token)
