local M ={}

local plugin_version = "0.0.1"
local function install_timepilot()
    if M.config.autoinstall then
        return os.execute("go install github.com/frenki123/timepilot.nvim/timepilot@latest")
    end
    return false
end

M.config = {
    timepilot_path = "timepilot",
    autoinstall = true
}


function M.set_config(user_config)
    for option, value in pairs(user_config) do
        M.config[option] = value
    end
end

function M.check_config()
    local timepilot = M.config.timepilot_path
    local handle = io.popen(timepilot .. "version")
    if not handle then
        if not install_timepilot() then
            print("Timepilot not installed")
        end
        return
    end
    local result = handle:read("*a")
    handle:close()
    local version = result:match("v(.+)")
    if version == plugin_version then
        return
    end
    if not install_timepilot() then
        print("Required Timepilot version is" .. plugin_version .. "installed" .. version)
    end
end

return M
