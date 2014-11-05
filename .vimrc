execute pathogen#infect()

"
" General settings
"
autocmd BufEnter *.json set filetype=javascript
filetype plugin indent on
let mapleader = ","
noremap <Right> <C-w>10<
noremap <Left> <C-w>10>
noremap <Up> <C-w>4+
noremap <Down> <C-w>4-

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
set listchars=tab:\▸-\,eol:¬,trail:·,nbsp:·
set tabstop=4
set shiftwidth=4
