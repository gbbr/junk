execute pathogen#infect()

"
" General settings
"
autocmd BufEnter *.json set filetype=javascript
filetype plugin indent on
let mapleader = ","

"
" Theme
"
syntax on
colorscheme apprentice
set cursorline
set number

"
" Autocomplete (SuperTab)
"
set completeopt=longest,menuone
let g:SuperTabDefaultCompletionType = "context"
let g:SuperTabContextDefaultCompletionType = "<c-n>"

"
" StatusBar (Airline)
"
let g:airline#extensions#tabline#enabled = 1
set laststatus=2

"
" EasyMotion
"
nmap s <Plug>(easymotion-s)

"
" NerdTree
" 
nmap <leader>n :NERDTreeToggle<CR>
nnoremap <Leader>f :NERDTreeFind<CR>
let NERDTreeQuitOnOpen=1

"
" Tabs and spaces
"
nmap <leader>l :set list!<CR>
set listchars=tab:▸\ ,eol:¬
set tabstop=4
