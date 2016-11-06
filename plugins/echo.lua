local plugin = require('plugin')

local echo = {
    Name = "echo",
    Desciption = "Echoes everything it hears",
    Bind = {
        OnMessage = function (bot, msg)
            bot:send(msg.Params[1], msg.Trailing)
        end
    }
}

plugin.add(echo)