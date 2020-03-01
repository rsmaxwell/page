
set BASE=..\..\output\html
xcopy /s /q /y html\* %BASE%\

set PREFIX=html/
set REQUEST_URI=http://example.com?image=pages/diary-1828-and-1829-and-jan-1830/img2834.jpg
pages.exe > %BASE%\page.html