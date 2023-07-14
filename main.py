from base_service import Service
from tortoise import Tortoise, run_async
from models import User, Session
from utils import SessionRediser, Rediser, verify_verification_token, create_verification_token
from argon2 import hash_password, verify_password
import uuid

async def on_user_password_create(data: dict):
    user = await User(
        password_hash= hash_password(
            data["password"].encode("utf-8")
        ).decode(),
        id=data["id"],
    )
    await user.save()
    vid = create_verification_token(user.id)
    return {"id": str(user.id), "vid": str(vid)}


async def on_user_password_update(data: dict):
    return {"abc": 1}

async def on_user_login(data: dict) -> dict:
    # TODO: implement ip check
    password: bytes = data["password"].encode("utf-8")
    user = await User.get(id=data["id"])
    logged_in: bool = verify_password(user.password_hash.encode("utf-8"), password)
    session = await Session.create(user=user.id, ip=str(data["ip"]))
    await (await session_redis.init()).new(session.id, user.id, 172800)
    #
    return {"id": str(session.id)}

async def on_user_session_exists(data: dict) -> bool:
    return await (await session_redis.init()).exists(data["sid"])


async def on_user_session_get_user(data: dict) -> str:
    return {await (await session_redis.init()).get(data["sid"])}

async def init():
    auth_service.logger.debug("connecting to database")
    await Tortoise.init(
        db_url="postgres://miauw_user:miauw_password@192.168.1.28:5432/miauw",
        modules={"models": ["models"]},
    )
    auth_service.logger.debug("generating schemas")
    await Tortoise.generate_schemas()


if __name__ == "__main__":
    session_redis = SessionRediser("redis://192.168.1.28")
    verfication_redis = Rediser("redis://192.168.1.28", 1)
    auth_service = Service("auth", "amqp://guest:guest@192.168.1.28")
    auth_service.add_event_handler("auth.password.initial", on_user_password_create)
    auth_service.add_event_handler("auth.password.change", on_user_password_update)
    auth_service.add_event_handler("auth.login", on_user_login)
    auth_service.add_event_handler("auth.session.exists", on_user_session_exists)
    auth_service.add_event_handler("auth.session.get-user", on_user_session_get_user)
    run_async(init())
    auth_service.start()
