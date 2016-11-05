local plugin = require('plugin')

local echo = {
    Name = "echo",
    Desciption = "Echoes everything it hears",
    Bind = {
        OnMessage = function (bot, msg)
            bot:send("#geoffrey", msg.Trailing)
        end
    }
}

plugin.add(echo)