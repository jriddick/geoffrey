local plugin = require('plugin')

local registration = {
    Name = "registration",
    Description = "Registers nick and user when connecting",
    Bind = {
        OnNotice = function (bot, msg)
            if msg.Trailing == "*** Looking up your hostname..." then
                bot:nick(bot.config.Nick)
                bot:user(bot.config.User, bot.config.Name)
            end
        end
    }
}

plugin.add(registration)