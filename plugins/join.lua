local plugin = require('plugin')

local join = {
    Name = "join",
    Description = "Join will join all pre-defined channels",
    Bind = {
        OnWelcome = function (bot, msg)
            for key, value in pairs(bot.config.Channels) do
                bot:join(value)
            end
        end
    }
}

plugin.add(join)