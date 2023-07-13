from tortoise import fields, Model, validators

class InetField(fields.Field):
    """
    IP Adress field specially for PostreSQL.
    This field can store an ip adress.
    """
    SQL_TYPE = "inet"

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

class Session(Model):
    id = fields.UUIDField(pk=True)
    user = fields.UUIDField()
    ip = fields.TextField()
    created_at = fields.DatetimeField(auto_now_add=True)

