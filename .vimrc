execute pathogen#infect()

" General settings
filetype plugin indent on
let mapleader = ","

" Source vimrc files after editing
autocmd bufwritepost .vimrc source <afile>

" Theme
syntax on
colorscheme apprentice
set cursorline
set number

" Autocomplete (SuperTab)
set completeopt=longest,menuone
let g:SuperTabDefaultCompletionType = "context"

" StatusBar (Airline)
let g:airline#extensions#tabline#enabled = 1
set laststatus=2

" EasyMotion
nmap s <Plug>(easymotion-s)

" NerdTree
nmap <leader>n :NERDTreeToggle<CR>
let NERDTreeQuitOnOpen=1

" Tabs and spaces
nmap <leader>l :set list!<CR>
set listchars=tab:▸\ ,eol:¬
set tabstop=4
