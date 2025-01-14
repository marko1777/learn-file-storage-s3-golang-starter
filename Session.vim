let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/projects/learn/learn-file-storage-s3-golang-starter
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +8 .env
badd +7 samplesdownload.sh
badd +133 main.go
badd +3 handler_get_thumbnail.go
badd +104 handler_upload_thumbnail.go
badd +129 internal/database/videos.go
badd +4 tubely.db
badd +1 reset.go
badd +31 assets.go
badd +1 assets/53ab1411-e398-4b1c-8063-05f0d0fdead6.png
badd +139 /usr/lib/go/src/mime/mediatype.go
badd +229 app/app.js
badd +1 cache.go
badd +133 /usr/lib/go/src/encoding/base64/base64.go
badd +197 handler_upload_video.go
badd +114 internal/auth/auth.go
badd +1 internal/database/refresh_tokens.go
badd +177 ~/go/pkg/mod/github.com/aws/aws-sdk-go-v2/service/s3@v1.72.0/api_op_PutObject.go
badd +45 ~/go/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.32.7/aws/to_ptr.go
badd +1 handler_users.go
badd +133 handler_video_meta.go
argglobal
%argdel
$argadd ~/projects/learn/learn-file-storage-s3-golang-starter/
edit handler_upload_video.go
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 94 + 48) / 96)
exe 'vert 2resize ' . ((&columns * 1 + 48) / 96)
argglobal
balt handler_video_meta.go
setlocal foldmethod=manual
setlocal foldexpr=0
setlocal foldmarker={{{,}}}
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldenable
silent! normal! zE
let &fdl = &fdl
let s:l = 1 - ((0 * winheight(0) + 23) / 46)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 1
normal! 0
lcd ~/projects/learn/learn-file-storage-s3-golang-starter
wincmd w
argglobal
if bufexists(fnamemodify("~/projects/learn/learn-file-storage-s3-golang-starter/.env", ":p")) | buffer ~/projects/learn/learn-file-storage-s3-golang-starter/.env | else | edit ~/projects/learn/learn-file-storage-s3-golang-starter/.env | endif
if &buftype ==# 'terminal'
  silent file ~/projects/learn/learn-file-storage-s3-golang-starter/.env
endif
setlocal foldmethod=manual
setlocal foldexpr=0
setlocal foldmarker={{{,}}}
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldenable
silent! normal! zE
let &fdl = &fdl
let s:l = 8 - ((7 * winheight(0) + 23) / 46)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 8
normal! 014|
lcd ~/projects/learn/learn-file-storage-s3-golang-starter
wincmd w
exe 'vert 1resize ' . ((&columns * 94 + 48) / 96)
exe 'vert 2resize ' . ((&columns * 1 + 48) / 96)
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
let g:this_session = v:this_session
let g:this_obsession = v:this_session
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
