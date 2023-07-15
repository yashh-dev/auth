from base_service import BaseService
from tortoise import Tortoise, run_async
from models import User, Session
from utils import (
    SessionRediser,
    Rediser,
    verify_verification_token,
    create_verification_token,
)
from argon2 import hash_password, verify_password
import uuid
import json
import jwt
import traceback


auth_service = BaseService("auth", "amqp://guest:guest@192.168.1.28")


@auth_service.event("auth.password.initial")
async def on_user_password_create(data: dict):
    user = await User(
        password_hash=hash_password(data["password"].encode("utf-8")).decode(),
        id=data["id"],
    )
    await user.save()
    return create_verification_token(user.id)


@auth_service.event("auth.password.change")
async def on_user_password_change(data: dict):
    return {"abc": 1}


@auth_service.event("auth.login")
async def on_user_login(data: dict) -> str:
    # TODO: implement ip check
    password: bytes = data["password"].encode("utf-8")
    user = await User.get(id=data["id"])
    logged_in: bool = verify_password(user.password_hash.encode("utf-8"), password)
    session = await Session.create(user=user.id, ip=str(data["ip"]))
    await (await session_redis.init()).new(session.id, user.id, 172800)
    #
    return str(session.id)


@auth_service.event("auth.session.exists")
async def on_user_session_exists(sid: str) -> bool:
    return await (await session_redis.init()).exists(sid)


@auth_service.event("auth.session.get-user")
async def on_user_session_get_user(sid: str) -> str:
    return await (await session_redis.init()).get(sid)


@auth_service.event("auth.verify")
async def on_user_verify(token: str) -> dict:
    try:
        payload = verify_verification_token(token)
        user = await User.get(id=payload["sub"])
        user.verified = True
        await user.save()
    except Exception as e:
        print(e)
        print(traceback.format_exc())
        return {"ok": 0}
    return {"ok": 1}


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
    run_async(init())
    auth_service.start()
