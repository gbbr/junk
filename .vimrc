execute pathogen#infect()

" General settings
syntax on
filetype plugin indent on
set number
set laststatus=2

" Airline
let g:airline#extensions#tabline#enabled = 1
let g:SuperTabDefaultCompletionType = "context"

" Keybindings
nmap <F3> :NERDTreeToggle<CR>
nmap s <Plug>(easymotion-s)
