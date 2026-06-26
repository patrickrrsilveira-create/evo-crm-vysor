import os
import glob
import re

def fix_success_flag():
    for filepath in glob.glob('src/api/*_routes.py'):
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Regex to find:
        # data={
        #    ... anything ...
        # },
        # message="Authorization completed successfully"
        
        # Actually, let's just replace `data={` with `data={"success": True, `
        # ONLY if the response is followed by `message="Authorization completed successfully"`
        
        # We can use a regex that matches `data=\{([^\}]*)\}(\s*,\s*message="Authorization completed successfully")`
        # and replaces it with `data={"success": True, \1}\2`
        
        pattern = r'data=\{([^\}]*)\}(\s*,\s*message="Authorization completed successfully")'
        
        # But what if there are nested brackets? e.g. data={"foo": result.get("email")} -> the inner bracket will not match [^\}]*
        # A simpler way: Find all occurrences of `message="Authorization completed successfully"`.
        # Find the preceding `data={`.
        
        # Let's split the file by `message="Authorization completed successfully"`
        parts = content.split('message="Authorization completed successfully"')
        
        if len(parts) > 1:
            new_content = parts[0]
            for i in range(1, len(parts)):
                part_before = new_content
                part_after = parts[i]
                
                # Find the last `data={` in part_before
                last_data_idx = part_before.rfind('data={')
                if last_data_idx != -1:
                    # Check if we already have "success": True
                    inner_str = part_before[last_data_idx:len(part_before)]
                    if '"success": True' not in inner_str and "'success': True" not in inner_str:
                        # Insert "success": True,
                        new_content = part_before[:last_data_idx + 6] + '"success": True, ' + part_before[last_data_idx + 6:]
                    else:
                        new_content = part_before
                else:
                    new_content = part_before
                    
                new_content += 'message="Authorization completed successfully"' + part_after
                
            if content != new_content:
                print(f"Fixed {filepath}")
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)

if __name__ == '__main__':
    fix_success_flag()
