local M = {}

function M.say_hello()
    vim.api.nvim_echo({{"Hello world", "Normal"}}, false, {})
end

return M
