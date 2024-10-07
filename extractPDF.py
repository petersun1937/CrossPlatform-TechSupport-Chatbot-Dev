import sys
import pdfplumber

def extract_text_from_pdf(file_path):
    extracted_text = ''
    with pdfplumber.open(file_path) as pdf:
        for page in pdf.pages:
            text = page.extract_text()
            if text:
                extracted_text += text + '\n'
    return extracted_text

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python pdf_extractor.py <path_to_pdf>")
        sys.exit(1)
    
    file_path = sys.argv[1]  # Get the file path from the command-line arguments
    try:
        text = extract_text_from_pdf(file_path)
        print(text)  # Print extracted text to stdout
    except Exception as e:
        print(f"Error extracting text: {e}")
        sys.exit(1)
