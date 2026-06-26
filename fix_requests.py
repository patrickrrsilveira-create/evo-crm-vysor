import os
import glob

def fix_files():
    for filepath in glob.glob('src/api/*.py'):
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        new_content = content.replace('request=request,', '').replace('request=request', '')
        
        if content != new_content:
            print(f"Fixed {filepath}")
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(new_content)

if __name__ == '__main__':
    fix_files()
