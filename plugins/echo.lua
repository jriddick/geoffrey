local plugin = require('plugin')

local echo = {
    Name = "echo",
    Desciption = "Echoes everything it hears",
    Bind = {
        OnMessage = function (bot, msg)
            if msg.Params[1] == bot.config.Nick then
                bot:send(msg.Prefix.Name, msg.Trailing)
            else
                bot:send(msg.Params[1], msg.Trailing)
            end
        end
    }
}

plugin.add(echo)