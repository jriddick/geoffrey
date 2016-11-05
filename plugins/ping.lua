local plugin = require('plugin')

local ping = {
    Name = "ping",
    Description = "Ping responds to pings",
    Bind = {
        OnPing = function (bot, msg)
            bot:pong(msg.Trailing)
        end
    }
}

plugin.add(ping)