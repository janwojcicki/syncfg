
set nocompatible              " be iMproved, required
filetype off                  " required

" set the runtime path to include Vundle and initialize
set rtp+=$HOME/.vim/bundle/Vundle.vim
call vundle#begin('$HOME/.vim/bundle/')
" alternatively, pass a path where Vundle should install plugins
"call vundle#begin('~/some/path/here')

" let Vundle manage Vundle, required
Plugin 'VundleVim/Vundle.vim'


" ==== helpers
Plugin 'vim-scripts/L9'

" ==== File tree
Plugin 'scrooloose/nerdtree'


Plugin 'vim-airline/vim-airline'
Plugin 'terryma/vim-expand-region'
Plugin 'flazz/vim-colorschemes'
Plugin 'valloric/youcompleteme'
Plugin 'ctrlpvim/ctrlp.vim'
Plugin 'tpope/vim-commentary'
Plugin 'vim-airline/vim-airline-themes'
Plugin 'raimondi/delimitmate'
Plugin 'vim-syntastic/syntastic'
Plugin 'rust-lang/rust.vim'
"Plugin 'w0rp/ale'
Plugin 'mileszs/ack.vim'
Plugin 'tpope/vim-fugitive'
Plugin 'farmergreg/vim-lastplace'
Plugin 'tomlion/vim-solidity'
Plugin 'rhysd/clever-f.vim'
Plugin 'mhartington/oceanic-next'
Plugin 'othree/yajs.vim'
Plugin 'rhysd/vim-crystal'
Plugin 'othree/html5.vim'
Plugin 'HerringtonDarkholme/yats.vim'
Plugin 'yuttie/comfortable-motion.vim'
Plugin 'lilydjwg/colorizer'
Plugin 'prettier/vim-prettier'
Plugin 'airblade/vim-gitgutter'
Plugin 'gryf/pylint-vim'
Plugin 'posva/vim-vue'
Plugin 'cakebaker/scss-syntax.vim'

call vundle#end()
filetype plugin indent on
let g:airline_theme='bubblegum'
set shiftwidth=4
set tabstop=4
set autoindent
set smartindent
map <C-t> :NERDTreeToggle<CR>
imap jk <Esc>

set background=dark
syntax on
map <C-J> <C-W>j<C-W>_
map <C-K> <C-W>k<C-W>_
set autochdir
set number
set encoding=utf-8


"easier split navigations
nnoremap <C-J> <C-W><C-J>
nnoremap <C-K> <C-W><C-K>
nnoremap <C-L> <C-W><C-L>
nnoremap <C-H> <C-W><C-H>

let g:ycm_global_ycm_extra_conf = '~/.ycm_extra_conf.py'
:set scrolloff=10
set wildignore+=*/tmp/*,*.so,*.swp,*.zip,*.o,*.d 

highlight clear CursorLine
highlight CursorLineNR ctermbg=67
highlight Cursor ctermbg=67
augroup CursorLine
    au!
    au VimEnter,WinEnter,BufWinEnter * setlocal cursorline
    au WinLeave * setlocal nocursorline
augroup END

set statusline+=%#warningmsg#
set statusline+=%{SyntasticStatuslineFlag()}
set statusline+=%*

let g:syntastic_always_populate_loc_list = 1
let g:syntastic_auto_loc_list = 1
let g:syntastic_check_on_open = 1
let g:syntastic_check_on_wq = 0

set mouse=a

if (has("termguicolors"))
 set termguicolors
endif

" Theme
syntax enable
colorscheme OceanicNext
let g:airline_theme = 'oceanicnext'

" " Copy to clipboard
vnoremap  Qy  "+y
nnoremap  QY  "+yg_
nnoremap  Qy  "+y
nnoremap  Qyy  "+yy

" " Paste from clipboard
nnoremap Qp "+p
nnoremap QP "+P
vnoremap Qp "+p
vnoremap QP "+P
set splitbelow
set splitright

let g:comfortable_motion_no_default_key_mappings = 1

nnoremap <silent> <C-s> :call comfortable_motion#flick(100)<CR>
nnoremap <silent> <C-a> :call comfortable_motion#flick(-100)<CR>

noremap <silent> <ScrollWheelDown> :call comfortable_motion#flick(40)<CR>
noremap <silent> <ScrollWheelUp>   :call comfortable_motion#flick(-40)<CR>
hi NonText guifg=bg

nnoremap <C-A-h> <C-w><char-60>
nnoremap <C-A-l> <C-w><char-62>

let g:prettier#autoformat = 0
autocmd BufWritePre *.js,*.jsx,*.mjs,*.ts,*.tsx,*.css,*.less,*.scss,*.json,*.graphql,*.md,*.vue Prettier

let g:NERDTreeWinSize=25

let g:prettier#config#tab_width = 4

let g:prettier#config#use_tabs = 'true'

let g:prettier#config#single_quote = 'true'

let g:ctrlp_custom_ignore = 'node_modules'

highlight OverLength ctermbg=red ctermfg=white guibg=#592929
match OverLength /\%99v.\+/
highlight Normal guibg=none
highlight NonText guibg=none
