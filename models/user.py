from tortoise import fields, Model

class User(Model):
    id = fields.UUIDField(pk=True)
    password_hash = fields.CharField(max_length=256)
    modified_at = fields.DatetimeField(auto_now=True)
    created_at = fields.DatetimeField(auto_now_add=True)

    class Meta:
        table = "users"