M = {}

local client = require("timepilot.client")

function M.autocmd()
  local augroup = vim.api.nvim_create_augroup("Timepilot", { clear = true })
  --Start client
  vim.api.nvim_create_autocmd("BufWinEnter", {
    group = augroup,
    once = true,
    callback = function()
      client.start()
    end,
  })

  --track when entered in insert mode
  vim.api.nvim_create_autocmd("InsertEnter", {
    callback = function()
      local filename = vim.api.nvim_buf_get_name(0)
      local basename = vim.fn.fnamemodify(filename, ":t")
      local params = {
        filetype = vim.bo.filetype,
        filename = basename,
        action = "InsertEnter",
      }
      client.send_event("action", params)
    end,
  })
end

return M
