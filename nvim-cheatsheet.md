Basic Navigation
h, j, k, l - Move cursor left, down, up, right (the core movement keys in Neovim)
w - Jump to start of next word (punctuation considered as words)
W - Jump to start of next WORD (space-separated words)
b - Jump to start of previous word
B - Jump to start of previous WORD
e - Jump to end of word
E - Jump to end of WORD
0 - Jump to start of line (first column)
^ - Jump to first non-blank character of line
$ - Jump to end of line
gg - Go to first line of document
G - Go to last line of document
{number}G - Go to specific line number
{ - Jump to previous paragraph/code block
} - Jump to next paragraph/code block
Ctrl-u - Move up half a screen
Ctrl-d - Move down half a screen
Ctrl-b - Move up one full screen
Ctrl-f - Move down one full screen
zz - Center cursor on screen (current line becomes middle line)
zt - Position cursor at top of screen
zb - Position cursor at bottom of screen
Editing Commands
i - Enter insert mode before cursor (for inserting text)
I - Enter insert mode at beginning of line
a - Enter insert mode after cursor (append)
A - Enter insert mode at end of line
o - Insert new line below current line and enter insert mode
O - Insert new line above current line and enter insert mode
r - Replace a single character under cursor (without entering insert mode)
R - Enter replace mode (overwriting existing text)
x - Delete character under cursor
X - Delete character before cursor
dd - Delete entire line (and store in register)
{number}dd - Delete multiple lines
D - Delete from cursor to end of line
yy or Y - Yank (copy) entire line
{number}yy - Yank multiple lines
y$ - Yank from cursor to end of line
p - Paste after cursor
P - Paste before cursor
u - Undo last change
Ctrl-r - Redo (undo the undo)
~ - Switch case of character under cursor
>> - Indent line
     << - Unindent line
     . - Repeat last command (powerful for repetitive edits)
     cc or C - Change entire line (delete line and enter insert mode)
     cw - Change word (delete word and enter insert mode)
     c$ or C - Change to end of line
     J - Join current line with the next line
     Search and Replace
     /pattern - Search forward for pattern
     ?pattern - Search backward for pattern
     n - Repeat search in same direction
     N - Repeat search in opposite direction
* - Search forward for word under cursor
# - Search backward for word under cursor
:%s/old/new/g - Replace all occurrences of 'old' with 'new' throughout file
:%s/old/new/gc - Replace all occurrences with confirmations
:s/old/new/g - Replace all occurrences on current line
:noh - Clear search highlighting
gd - Go to local definition of word under cursor
gD - Go to global definition of word under cursor
Visual Mode
v - Enter character-wise visual mode (select characters)
V - Enter line-wise visual mode (select entire lines)
Ctrl-v - Enter block-wise visual mode (select rectangular blocks)
gv - Reselect previous visual selection
o - In visual mode: Move to other end of selection
O - In visual block mode: Move to other corner of block
aw - Select a word (in visual mode)
ab - Select a block with () (in visual mode)
aB - Select a block with {} (in visual mode)
at - Select a block with HTML/XML tags (in visual mode)
File Operations
:e filename - Edit a file (create if doesn't exist)
:w - Write (save) the file
:w filename - Write to specified filename (save as)
:q - Quit (fails if unsaved changes)
:q! - Quit without saving (discard changes)
:wq or :x - Write and quit
:saveas filename - Save file as filename
:r filename - Insert contents of file below cursor
:r !command - Insert output of shell command below cursor
Working with Windows and Tabs
:split or :sp - Split window horizontally
:vsplit or :vs - Split window vertically
Ctrl-w h/j/k/l - Navigate between windows (left/down/up/right)
Ctrl-w +/- - Increase/decrease window height
Ctrl-w </>- Increase/decrease window width
Ctrl-w = - Make all windows equal size
Ctrl-w o - Make current window the only one
:tabnew - Create new tab
gt - Go to next tab
gT - Go to previous tab
:tabclose - Close current tab
:tabonly - Close all other tabs
Buffer Management
:ls - List all buffers
:b number - Switch to buffer by number
:bn - Next buffer
:bp - Previous buffer
:bd - Delete buffer (close file)
:bufdo command - Execute command on all buffers
:e # - Edit the alternate file (usually the previously edited file)
Marks and Jumps
m{a-z} - Set mark at current position (lowercase for file-local)
m{A-Z} - Set mark at current position (uppercase for global)
'{mark} - Jump to line of mark
`{mark} - Jump to position of mark
Ctrl-o - Jump to older position in jump list
Ctrl-i - Jump to newer position in jump list
'. - Jump to position of last change
`. - Jump to exact position of last change
Text Objects and Motions
ci( - Change inside parentheses
di" - Delete inside double quotes
yi] - Yank inside square brackets
va{ - Visually select around curly braces (including the braces)
dap - Delete around paragraph
cit - Change inside HTML/XML tag
diw - Delete inside word
daw - Delete around word (including spaces)
dab - Delete around block (parentheses)
daB - Delete around block (curly braces)
Fold Commands
zf - Create fold (in visual mode)
zo - Open fold under cursor
zc - Close fold under cursor
za - Toggle fold under cursor
zR - Open all folds
zM - Close all folds
zj - Move to next fold
zk - Move to previous fold
Neovim-Specific Features
:terminal or :term - Open integrated terminal
Ctrl-\ Ctrl-n - Exit terminal mode to normal mode
:checkhealth - Run Neovim's diagnostic tool
:lua require('telescope.builtin').find_files() - Use Telescope plugin to find files
:TSInstall language - Install treesitter parser for a language
:LspInfo - Show Language Server Protocol status
:TSBufToggle highlight - Toggle treesitter highlighting
:highlight - Show current highlight groups
:Tutor - Start Neovim's built-in tutorial
:help nvim-features - View Neovim's specific features
Advanced Features
q{a-z} - Record macro into register
@{a-z} - Play macro from register
@@ - Repeat last played macro
g& - Repeat last substitution on all lines
:norm cmd - Execute normal mode command on selected lines
gf - Go to file under cursor
Ctrl-a - Increment number under cursor
Ctrl-x - Decrement number under cursor
:sort - Sort selected lines
!motion command - Filter text through external command




----------------------
```bash
diskutil list
diskutil eraseDisk MS-DOS "WIN11" MBR /dev/disk5
rsync -avh --progress --exclude=sources/install.wim /Volumes/CCCOMA_X64FRE_EN-US_DV9/ /Volumes/WIN11
brew install wimlib
wimlib-imagex split /Volumes/CCCOMA_X64FRE_EN-US_DV9/sources/install.wim /Volumes/WIN11/sources/install.swm 3800

```
