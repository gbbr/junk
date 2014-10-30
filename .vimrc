execute pathogen#infect()

" General settings
syntax on
filetype plugin indent on

set cursorline
set number
set laststatus=2
set tabstop=4
set listchars=tab:▸\ ,eol:¬

let NERDTreeQuitOnOpen=1
let mapleader = ","

" Colorscheme
colorscheme apprentice

" Autocomplete
set completeopt=longest,menuone
let g:SuperTabDefaultCompletionType = "context"

" Airline
let g:airline#extensions#tabline#enabled = 1

" Keybindings
nmap <leader>n :NERDTreeToggle<CR>
nmap s <Plug>(easymotion-s)
nmap <leader>l :set list!<CR>
