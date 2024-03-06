--- TODO: This is a copy of the system.lua file from the neovim source code.
--- https://github.com/neovim/neovim/blob/master/runtime/lua/vim/_system.lua
---
--- If any coredev sees this and feels i am breaking any license, please let me
--- know and i'll move over to plenary, just didn't want to legacy apis
--
-- LICENSE
--Copyright Neovim contributors. All rights reserved.
--
-- Neovim is licensed under the terms of the Apache 2.0 license, except for
-- parts of Neovim that were contributed under the Vim license (see below).
--
-- Neovim's license follows:
--
-- ====
--                                  Apache License
--                            Version 2.0, January 2004
--                         https://www.apache.org/licenses/
--
--    TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
--
--    1. Definitions.
--
--       "License" shall mean the terms and conditions for use, reproduction,
--       and distribution as defined by Sections 1 through 9 of this document.
--
--       "Licensor" shall mean the copyright owner or entity authorized by
--       the copyright owner that is granting the License.
--
--       "Legal Entity" shall mean the union of the acting entity and all
--       other entities that control, are controlled by, or are under common
--       control with that entity. For the purposes of this definition,
--       "control" means (i) the power, direct or indirect, to cause the
--       direction or management of such entity, whether by contract or
--       otherwise, or (ii) ownership of fifty percent (50%) or more of the
--       outstanding shares, or (iii) beneficial ownership of such entity.
--
--       "You" (or "Your") shall mean an individual or Legal Entity
--       exercising permissions granted by this License.
--
--       "Source" form shall mean the preferred form for making modifications,
--       including but not limited to software source code, documentation
--       source, and configuration files.
--
--       "Object" form shall mean any form resulting from mechanical
--       transformation or translation of a Source form, including but
--       not limited to compiled object code, generated documentation,
--       and conversions to other media types.
--
--       "Work" shall mean the work of authorship, whether in Source or
--       Object form, made available under the License, as indicated by a
--       copyright notice that is included in or attached to the work
--       (an example is provided in the Appendix below).
--
--       "Derivative Works" shall mean any work, whether in Source or Object
--       form, that is based on (or derived from) the Work and for which the
--       editorial revisions, annotations, elaborations, or other modifications
--       represent, as a whole, an original work of authorship. For the purposes
--       of this License, Derivative Works shall not include works that remain
--       separable from, or merely link (or bind by name) to the interfaces of,
--       the Work and Derivative Works thereof.
--
--       "Contribution" shall mean any work of authorship, including
--       the original version of the Work and any modifications or additions
--       to that Work or Derivative Works thereof, that is intentionally
--       submitted to Licensor for inclusion in the Work by the copyright owner
--       or by an individual or Legal Entity authorized to submit on behalf of
--       the copyright owner. For the purposes of this definition, "submitted"
--       means any form of electronic, verbal, or written communication sent
--       to the Licensor or its representatives, including but not limited to
--       communication on electronic mailing lists, source code control systems,
--       and issue tracking systems that are managed by, or on behalf of, the
--       Licensor for the purpose of discussing and improving the Work, but
--       excluding communication that is conspicuously marked or otherwise
--       designated in writing by the copyright owner as "Not a Contribution."
--
--       "Contributor" shall mean Licensor and any individual or Legal Entity
--       on behalf of whom a Contribution has been received by Licensor and
--       subsequently incorporated within the Work.
--
--    2. Grant of Copyright License. Subject to the terms and conditions of
--       this License, each Contributor hereby grants to You a perpetual,
--       worldwide, non-exclusive, no-charge, royalty-free, irrevocable
--       copyright license to reproduce, prepare Derivative Works of,
--       publicly display, publicly perform, sublicense, and distribute the
--       Work and such Derivative Works in Source or Object form.
--
--    3. Grant of Patent License. Subject to the terms and conditions of
--       this License, each Contributor hereby grants to You a perpetual,
--       worldwide, non-exclusive, no-charge, royalty-free, irrevocable
--       (except as stated in this section) patent license to make, have made,
--       use, offer to sell, sell, import, and otherwise transfer the Work,
--       where such license applies only to those patent claims licensable
--       by such Contributor that are necessarily infringed by their
--       Contribution(s) alone or by combination of their Contribution(s)
--       with the Work to which such Contribution(s) was submitted. If You
--       institute patent litigation against any entity (including a
--       cross-claim or counterclaim in a lawsuit) alleging that the Work
--       or a Contribution incorporated within the Work constitutes direct
--       or contributory patent infringement, then any patent licenses
--       granted to You under this License for that Work shall terminate
--       as of the date such litigation is filed.
--
--    4. Redistribution. You may reproduce and distribute copies of the
--       Work or Derivative Works thereof in any medium, with or without
--       modifications, and in Source or Object form, provided that You
--       meet the following conditions:
--
--       (a) You must give any other recipients of the Work or
--           Derivative Works a copy of this License; and
--
--       (b) You must cause any modified files to carry prominent notices
--           stating that You changed the files; and
--
--       (c) You must retain, in the Source form of any Derivative Works
--           that You distribute, all copyright, patent, trademark, and
--           attribution notices from the Source form of the Work,
--           excluding those notices that do not pertain to any part of
--           the Derivative Works; and
--
--       (d) If the Work includes a "NOTICE" text file as part of its
--           distribution, then any Derivative Works that You distribute must
--           include a readable copy of the attribution notices contained
--           within such NOTICE file, excluding those notices that do not
--           pertain to any part of the Derivative Works, in at least one
--           of the following places: within a NOTICE text file distributed
--           as part of the Derivative Works; within the Source form or
--           documentation, if provided along with the Derivative Works; or,
--           within a display generated by the Derivative Works, if and
--           wherever such third-party notices normally appear. The contents
--           of the NOTICE file are for informational purposes only and
--           do not modify the License. You may add Your own attribution
--           notices within Derivative Works that You distribute, alongside
--           or as an addendum to the NOTICE text from the Work, provided
--           that such additional attribution notices cannot be construed
--           as modifying the License.
--
--       You may add Your own copyright statement to Your modifications and
--       may provide additional or different license terms and conditions
--       for use, reproduction, or distribution of Your modifications, or
--       for any such Derivative Works as a whole, provided Your use,
--       reproduction, and distribution of the Work otherwise complies with
--       the conditions stated in this License.
--
--    5. Submission of Contributions. Unless You explicitly state otherwise,
--       any Contribution intentionally submitted for inclusion in the Work
--       by You to the Licensor shall be under the terms and conditions of
--       this License, without any additional terms or conditions.
--       Notwithstanding the above, nothing herein shall supersede or modify
--       the terms of any separate license agreement you may have executed
--       with Licensor regarding such Contributions.
--
--    6. Trademarks. This License does not grant permission to use the trade
--       names, trademarks, service marks, or product names of the Licensor,
--       except as required for reasonable and customary use in describing the
--       origin of the Work and reproducing the content of the NOTICE file.
--
--    7. Disclaimer of Warranty. Unless required by applicable law or
--       agreed to in writing, Licensor provides the Work (and each
--       Contributor provides its Contributions) on an "AS IS" BASIS,
--       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
--       implied, including, without limitation, any warranties or conditions
--       of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
--       PARTICULAR PURPOSE. You are solely responsible for determining the
--       appropriateness of using or redistributing the Work and assume any
--       risks associated with Your exercise of permissions under this License.
--
--    8. Limitation of Liability. In no event and under no legal theory,
--       whether in tort (including negligence), contract, or otherwise,
--       unless required by applicable law (such as deliberate and grossly
--       negligent acts) or agreed to in writing, shall any Contributor be
--       liable to You for damages, including any direct, indirect, special,
--       incidental, or consequential damages of any character arising as a
--       result of this License or out of the use or inability to use the
--       Work (including but not limited to damages for loss of goodwill,
--       work stoppage, computer failure or malfunction, or any and all
--       other commercial damages or losses), even if such Contributor
--       has been advised of the possibility of such damages.
--
--    9. Accepting Warranty or Additional Liability. While redistributing
--       the Work or Derivative Works thereof, You may choose to offer,
--       and charge a fee for, acceptance of support, warranty, indemnity,
--       or other liability obligations and/or rights consistent with this
--       License. However, in accepting such obligations, You may act only
--       on Your own behalf and on Your sole responsibility, not on behalf
--       of any other Contributor, and only if You agree to indemnify,
--       defend, and hold each Contributor harmless for any liability
--       incurred by, or claims asserted against, such Contributor by reason
--       of your accepting any such warranty or additional liability.
--
-- ====
--
-- The above license applies to all parts of Neovim except (1) parts that were
-- contributed under the Vim license and (2) externally maintained libraries.
--
-- The externally maintained libraries used by Neovim are:
--
--   - Klib: a Generic Library in C. MIT/X11 license.
--   - Lua: MIT license
--   - LuaJIT: a Just-In-Time Compiler for Lua. Copyright Mike Pall. MIT license.
--   - Luv: Apache 2.0 license
--   - libmpack: MIT license
--   - libtermkey: MIT license
--   - libuv. Copyright Joyent, Inc. and other Node contributors. Node.js license.
--   - libvterm: MIT license
--   - lua-cjson: MIT license
--   - lua-compat: MIT license
--   - tree-sitter: MIT license
--   - unibilium: LGPL v3
--   - xdiff: LGPL v2
--
-- ====
--
-- Any parts of Neovim that were contributed under the Vim license are licensed
-- under the Vim license unless the copyright holder gave permission to license
-- those contributions under the Apache 2.0 license.
--
-- The Vim license follows:
--
-- VIM LICENSE
--
-- I)  There are no restrictions on distributing unmodified copies of Vim except
--     that they must include this license text.  You can also distribute
--     unmodified parts of Vim, likewise unrestricted except that they must
--     include this license text.  You are also allowed to include executables
--     that you made from the unmodified Vim sources, plus your own usage
--     examples and Vim scripts.
--
-- II) It is allowed to distribute a modified (or extended) version of Vim,
--     including executables and/or source code, when the following four
--     conditions are met:
--     1) This license text must be included unmodified.
--     2) The modified Vim must be distributed in one of the following five ways:
--        a) If you make changes to Vim yourself, you must clearly describe in
-- 	  the distribution how to contact you.  When the maintainer asks you
-- 	  (in any way) for a copy of the modified Vim you distributed, you
-- 	  must make your changes, including source code, available to the
-- 	  maintainer without fee.  The maintainer reserves the right to
-- 	  include your changes in the official version of Vim.  What the
-- 	  maintainer will do with your changes and under what license they
-- 	  will be distributed is negotiable.  If there has been no negotiation
-- 	  then this license, or a later version, also applies to your changes.
-- 	  The current maintainers are listed here: https://github.com/orgs/vim/people.
-- 	  If this changes it will be announced in appropriate places (most likely
-- 	  vim.sf.net, www.vim.org and/or comp.editors).  When it is completely
-- 	  impossible to contact the maintainer, the obligation to send him
-- 	  your changes ceases.  Once the maintainer has confirmed that he has
-- 	  received your changes they will not have to be sent again.
--        b) If you have received a modified Vim that was distributed as
-- 	  mentioned under a) you are allowed to further distribute it
-- 	  unmodified, as mentioned at I).  If you make additional changes the
-- 	  text under a) applies to those changes.
--        c) Provide all the changes, including source code, with every copy of
-- 	  the modified Vim you distribute.  This may be done in the form of a
-- 	  context diff.  You can choose what license to use for new code you
-- 	  add.  The changes and their license must not restrict others from
-- 	  making their own changes to the official version of Vim.
--        d) When you have a modified Vim which includes changes as mentioned
-- 	  under c), you can distribute it without the source code for the
-- 	  changes if the following three conditions are met:
-- 	  - The license that applies to the changes permits you to distribute
-- 	    the changes to the Vim maintainer without fee or restriction, and
-- 	    permits the Vim maintainer to include the changes in the official
-- 	    version of Vim without fee or restriction.
-- 	  - You keep the changes for at least three years after last
-- 	    distributing the corresponding modified Vim.  When the maintainer
-- 	    or someone who you distributed the modified Vim to asks you (in
-- 	    any way) for the changes within this period, you must make them
-- 	    available to him.
-- 	  - You clearly describe in the distribution how to contact you.  This
-- 	    contact information must remain valid for at least three years
-- 	    after last distributing the corresponding modified Vim, or as long
-- 	    as possible.
--        e) When the GNU General Public License (GPL) applies to the changes,
-- 	  you can distribute the modified Vim under the GNU GPL version 2 or
-- 	  any later version.
--     3) A message must be added, at least in the output of the ":version"
--        command and in the intro screen, such that the user of the modified Vim
--        is able to see that it was modified.  When distributing as mentioned
--        under 2)e) adding the message is only required for as far as this does
--        not conflict with the license used for the changes.
--     4) The contact information as required under 2)a) and 2)d) must not be
--        removed or changed, except that the person himself can make
--        corrections.
--
-- III) If you distribute a modified version of Vim, you are encouraged to use
--      the Vim license for your changes and make them available to the
--      maintainer, including the source code.  The preferred way to do this is
--      by e-mail or by uploading the files to a server and e-mailing the URL.
--      If the number of changes is small (e.g., a modified Makefile) e-mailing a
--      context diff will do.  The e-mail address to be used is
--      <maintainer@vim.org>
--
-- IV)  It is not allowed to remove this license from the distribution of the Vim
--      sources, parts of it or from a modified version.  You may use this
--      license for previous Vim releases instead of the license that they came
--      with, at your option.


local uv = vim.loop

--- @class vim.SystemOpts
--- @field stdin? string|string[]|true
--- @field stdout? fun(err:string?, data: string?)|false
--- @field stderr? fun(err:string?, data: string?)|false
--- @field cwd? string
--- @field env? table<string,string|number>
--- @field clear_env? boolean
--- @field text? boolean
--- @field timeout? integer Timeout in ms
--- @field detach? boolean

--- @class vim.SystemCompleted
--- @field code integer
--- @field signal integer
--- @field stdout? string
--- @field stderr? string

--- @class vim.SystemState
--- @field handle? uv.uv_process_t
--- @field timer?  uv.uv_timer_t
--- @field pid? integer
--- @field timeout? integer
--- @field done? boolean|'timeout'
--- @field stdin? uv.uv_stream_t
--- @field stdout? uv.uv_stream_t
--- @field stderr? uv.uv_stream_t
--- @field stdout_data? string[]
--- @field stderr_data? string[]
--- @field result? vim.SystemCompleted

--- @enum vim.SystemSig
local SIG = {
  HUP = 1, -- Hangup
  INT = 2, -- Interrupt from keyboard
  KILL = 9, -- Kill signal
  TERM = 15, -- Termination signal
  -- STOP = 17,19,23  -- Stop the process
}

---@param handle uv.uv_handle_t?
local function close_handle(handle)
  if handle and not handle:is_closing() then
    handle:close()
  end
end

---@param state vim.SystemState
local function close_handles(state)
  close_handle(state.handle)
  close_handle(state.stdin)
  close_handle(state.stdout)
  close_handle(state.stderr)
  close_handle(state.timer)
end

--- @class vim.SystemObj
--- @field pid integer
--- @field private _state vim.SystemState
--- @field wait fun(self: vim.SystemObj, timeout?: integer): vim.SystemCompleted
--- @field kill fun(self: vim.SystemObj, signal: integer|string)
--- @field write fun(self: vim.SystemObj, data?: string|string[])
--- @field is_closing fun(self: vim.SystemObj): boolean
local SystemObj = {}

--- @param state vim.SystemState
--- @return vim.SystemObj
local function new_systemobj(state)
  return setmetatable({
    pid = state.pid,
    _state = state,
  }, { __index = SystemObj })
end

--- @param signal integer|string
function SystemObj:kill(signal)
  self._state.handle:kill(signal)
end

--- @package
--- @param signal? vim.SystemSig
function SystemObj:_timeout(signal)
  self._state.done = 'timeout'
  self:kill(signal or SIG.TERM)
end

local MAX_TIMEOUT = 2 ^ 31

--- @param timeout? integer
--- @return vim.SystemCompleted
function SystemObj:wait(timeout)
  local state = self._state

  local done = vim.wait(timeout or state.timeout or MAX_TIMEOUT, function()
    return state.result ~= nil
  end, nil, true)

  if not done then
    -- Send sigkill since this cannot be caught
    self:_timeout(SIG.KILL)
    vim.wait(timeout or state.timeout or MAX_TIMEOUT, function()
      return state.result ~= nil
    end, nil, true)
  end

  return state.result
end

--- @param data string[]|string|nil
function SystemObj:write(data)
  local stdin = self._state.stdin

  if not stdin then
    error('stdin has not been opened on this object')
  end

  if type(data) == 'table' then
    for _, v in ipairs(data) do
      stdin:write(v)
      stdin:write('\n')
    end
  elseif type(data) == 'string' then
    stdin:write(data)
  elseif data == nil then
    -- Shutdown the write side of the duplex stream and then close the pipe.
    -- Note shutdown will wait for all the pending write requests to complete
    -- TODO(lewis6991): apparently shutdown doesn't behave this way.
    -- (https://github.com/neovim/neovim/pull/17620#discussion_r820775616)
    stdin:write('', function()
      stdin:shutdown(function()
        if stdin then
          stdin:close()
        end
      end)
    end)
  end
end

--- @return boolean
function SystemObj:is_closing()
  local handle = self._state.handle
  return handle == nil or handle:is_closing() or false
end

---@param output fun(err:string?, data: string?)|false
---@return uv.uv_stream_t?
---@return fun(err:string?, data: string?)? Handler
local function setup_output(output)
  if output == nil then
    return assert(uv.new_pipe(false)), nil
  end

  if type(output) == 'function' then
    return assert(uv.new_pipe(false)), output
  end

  assert(output == false)
  return nil, nil
end

---@param input string|string[]|true|nil
---@return uv.uv_stream_t?
---@return string|string[]?
local function setup_input(input)
  if not input then
    return
  end

  local towrite --- @type string|string[]?
  if type(input) == 'string' or type(input) == 'table' then
    towrite = input
  end

  return assert(uv.new_pipe(false)), towrite
end

--- @return table<string,string>
local function base_env()
  local env = vim.fn.environ() --- @type table<string,string>
  env['NVIM'] = vim.v.servername
  env['NVIM_LISTEN_ADDRESS'] = nil
  return env
end

--- uv.spawn will completely overwrite the environment
--- when we just want to modify the existing one, so
--- make sure to prepopulate it with the current env.
--- @param env? table<string,string|number>
--- @param clear_env? boolean
--- @return string[]?
local function setup_env(env, clear_env)
  if clear_env then
    return env
  end

  --- @type table<string,string|number>
  env = vim.tbl_extend('force', base_env(), env or {})

  local renv = {} --- @type string[]
  for k, v in pairs(env) do
    renv[#renv + 1] = string.format('%s=%s', k, tostring(v))
  end

  return renv
end

--- @param stream uv.uv_stream_t
--- @param text? boolean
--- @param bucket string[]
--- @return fun(err: string?, data: string?)
local function default_handler(stream, text, bucket)
  return function(err, data)
    if err then
      error(err)
    end
    if data ~= nil then
      if text then
        bucket[#bucket + 1] = data:gsub('\r\n', '\n')
      else
        bucket[#bucket + 1] = data
      end
    else
      stream:read_stop()
      stream:close()
    end
  end
end

local M = {}

--- @param cmd string
--- @param opts uv.spawn.options
--- @param on_exit fun(code: integer, signal: integer)
--- @param on_error fun()
--- @return uv.uv_process_t, integer
local function spawn(cmd, opts, on_exit, on_error)
  local handle, pid_or_err = uv.spawn(cmd, opts, on_exit)
  if not handle then
    on_error()
    error(pid_or_err)
  end
  return handle, pid_or_err --[[@as integer]]
end

---@param timeout integer
---@param cb fun()
---@return uv.uv_timer_t
local function timer_oneshot(timeout, cb)
  local timer = assert(uv.new_timer())
  timer:start(timeout, 0, function()
    timer:stop()
    timer:close()
    cb()
  end)
  return timer
end

--- @param state vim.SystemState
--- @param code integer
--- @param signal integer
--- @param on_exit fun(result: vim.SystemCompleted)?
local function _on_exit(state, code, signal, on_exit)
  close_handles(state)

  local check = assert(uv.new_check())
  check:start(function()
    for _, pipe in pairs({ state.stdin, state.stdout, state.stderr }) do
      if not pipe:is_closing() then
        return
      end
    end
    check:stop()
    check:close()

    if state.done == nil then
      state.done = true
    end

    if (code == 0 or code == 1) and state.done == 'timeout' then
      -- Unix: code == 0
      -- Windows: code == 1
      code = 124
    end

    local stdout_data = state.stdout_data
    local stderr_data = state.stderr_data

    state.result = {
      code = code,
      signal = signal,
      stdout = stdout_data and table.concat(stdout_data) or nil,
      stderr = stderr_data and table.concat(stderr_data) or nil,
    }

    if on_exit then
      on_exit(state.result)
    end
  end)
end

--- Run a system command
---
--- @param cmd string[]
--- @param opts? vim.SystemOpts
--- @param on_exit? fun(out: vim.SystemCompleted)
--- @return vim.SystemObj
function M.run(cmd, opts, on_exit)
  vim.validate({
    cmd = { cmd, 'table' },
    opts = { opts, 'table', true },
    on_exit = { on_exit, 'function', true },
  })

  opts = opts or {}

  local stdout, stdout_handler = setup_output(opts.stdout)
  local stderr, stderr_handler = setup_output(opts.stderr)
  local stdin, towrite = setup_input(opts.stdin)

  --- @type vim.SystemState
  local state = {
    done = false,
    cmd = cmd,
    timeout = opts.timeout,
    stdin = stdin,
    stdout = stdout,
    stderr = stderr,
  }

  --- @diagnostic disable-next-line:missing-fields
  state.handle, state.pid = spawn(cmd[1], {
    args = vim.list_slice(cmd, 2),
    stdio = { stdin, stdout, stderr },
    cwd = opts.cwd,
    --- @diagnostic disable-next-line:assign-type-mismatch
    env = setup_env(opts.env, opts.clear_env),
    detached = opts.detach,
    hide = true,
  }, function(code, signal)
    _on_exit(state, code, signal, on_exit)
  end, function()
    close_handles(state)
  end)

  if stdout then
    state.stdout_data = {}
    stdout:read_start(stdout_handler or default_handler(stdout, opts.text, state.stdout_data))
  end

  if stderr then
    state.stderr_data = {}
    stderr:read_start(stderr_handler or default_handler(stderr, opts.text, state.stderr_data))
  end

  local obj = new_systemobj(state)

  if towrite then
    obj:write(towrite)
    obj:write(nil) -- close the stream
  end

  if opts.timeout then
    state.timer = timer_oneshot(opts.timeout, function()
      if state.handle and state.handle:is_active() then
        obj:_timeout()
      end
    end)
  end

  return obj
end

return M

