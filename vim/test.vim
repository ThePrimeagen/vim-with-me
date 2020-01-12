imap <buffer> <silent> <expr> <F12> Double("\<F12>")
function! Double(mymap)
  try
    let char = getchar()
  catch /^Vim:Interrupt$/
    let char = "\<Esc>"
  endtry
  "exec BPBreakIf(char == 32, 1)
  if char == '^\d\+$' || type(char) == 0
    let char = nr2char(char)
  endif " It is the ascii code.
  if char == "\<Esc>"
    return ''
  endif
  redraw
  return char.char."\<C-R>=Redraw()\<CR>".a:mymap
endfunction

function! Redraw()
  redraw
  return ''
endfunction

