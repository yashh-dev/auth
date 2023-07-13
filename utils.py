from redis.asyncio import BlockingConnectionPool, Redis
import uuid

class Rediser:
    def __init__(self, url: str) -> None:
        self.url = url
        self.__pool = BlockingConnectionPool(host="192.168.1.28")
        self.radis: Redis = None

    async def init(self) -> "Rediser":
        self.radis = Redis(host="192.168.1.28", connection_pool=self.__pool, auto_close_connection_pool=True)
        return self

    async def set(key: str, value: str) -> None:
        await self.radis.set(key, value)

    async def get(self, key: str):
        """gets value from key"""
        return await self.radis.get(key)


class SessionRediser(Rediser):
    def __init__(self, url: str):
        super().__init__(url)

    async def exists(self, sid: str) -> bool:
        """checks if session the user is trying to use exists"""
        uid = await self.get(sid)
        if not uid:
            # session doesnt exist
            return False
        return True

    async def new(self, sid: uuid.UUID, uid: uuid.UUID, ttl: int):
        await self.radis.set(str(sid), str(uid), ex=172800)

