import requests
import sys

API_ENTRY_POINT1 = "https://api.ilovepdf.com/v1/auth"
API_ENTRY_POINT2="https://api.ilovepdf.com/v1/start"


class ILovePdf:
    def __init__(self,public_key):
        self.public_key=public_key
        self.token=self.get_token()
        
    def get_token(self):

        data={
            "public_key":self.public_key
        }
        response=requests.post(API_ENTRY_POINT1,data).json()
        return response['token']

class operations(ILovePdf):
    
    def __init__(self, public_key):
        super().__init__(public_key)
        self.headers = {"Authorization": "Bearer {}".format(self.token)}
        
    
    def start_task(self,tool):
        self.files = []
        self.tool = tool
        url = "{}/{}".format(API_ENTRY_POINT2,self.tool)
        response=requests.get(url,headers=self.headers).json()
        self.server = response["server"]
        self.task_id = response["task"]
        self.base_api_url = "https://{}/v1".format(self.server)
        return response
    
    def add_file(self, filename):
        url = self.base_api_url + "/upload"
        params = {"task": self.task_id}
        files = {"file": open(filename, "rb")}
        response = requests.post(
            url,
            params,
            files=files,
            headers=self.headers
        ).json()
        self.server_filename = response["server_filename"]
        self.files.append({
            "server_filename": response["server_filename"],
            "filename": filename
        })
        return response
    
    def execute(self, password):
        url = self.base_api_url + "/process"
        fixed_params = {
            "task": self.task_id,
            "tool": self.tool,
            "files": self.files,
            "password" : password
        }
        params = fixed_params.copy()
        response = requests.post(url, json=params, headers=self.headers).json()
        self.timer = response["timer"]
        return response
    
    def download(self, output_filename):
        url = self.base_api_url + "/download/{}".format(self.task_id)
        response = requests.get(url, headers=self.headers)
        with open(output_filename, "wb") as output_file:
            output_file.write(response.content)
        return response.status_code    



if __name__ == "__main__":

    public_key = "project_public_db7deec963dc9219b319768d2766bfc6_9-1mScb0a712112737d004c62656bb16f2eb1"
    i=operations(public_key)

    while True:
        i.start_task("protect")

        path=input("Escriba la ruta donde se encuentra el archivo (incluya el nombre): ")
        contra = input("Escriba el ID del cliente objetivo: ")

        i.add_file(path.replace('"', ''))
        i.execute(contra)
        file_name = path.split("\\")
        file_name = path.split(".")
        file_name = file_name[0] + "_protegido.pdf"
        i.download(file_name)
        
        if input("Do you want to quit? (y/n) : ").lower() == "y":
            sys.exit()
