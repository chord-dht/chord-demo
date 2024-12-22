import os.path
import shutil
import subprocess
import time

from utils import generate_hash
from utils import thread_print


class Node:

    def __init__(self, work_dir, process_name, bind_host, bind_port):
        self.work_dir = work_dir
        self.process_name = process_name
        self.bind_host = bind_host
        self.bind_port = bind_port
        self.process = None
        self.identifier = generate_hash(bind_host + ":" + bind_port)
        if os.path.exists(self.work_dir):
            shutil.rmtree(self.work_dir)
        os.makedirs(self.work_dir)
        shutil.copy("./chord", self.work_dir)

    def create(self):
        params = [
            './chord',
            '-a', self.bind_host,
            '-p', self.bind_port,
            '--ts', '3000',
            '--tff', '1000',
            '--tcp', '3000',
            '-r', '4'
        ]
        thread_print("Starting Go process : ", self.process_name)
        output_file_path = os.path.join(self.work_dir, self.process_name + "_output.log")
        with open(output_file_path, "w") as output_file:
            process = subprocess.Popen(
                params,
                cwd=self.work_dir,
                stdin=subprocess.PIPE,
                stdout=output_file,
                stderr=subprocess.PIPE,
                text=True
            )
            self.process = process
            return process

    def join(self, join_node):

        params = [
            './chord',
            '-a', self.bind_host,
            '-p', self.bind_port,
            '-ja', join_node.bind_host,
            '-jp', join_node.bind_port,
            '--ts', '3000',
            '--tff', '1000',
            '--tcp', '3000',
            '-r', '4'
        ]
        thread_print("Starting Go process : ", self.process_name)
        output_file_path = os.path.join(self.work_dir, self.process_name + "_output.log")
        with open(output_file_path, "w") as output_file:
            process = subprocess.Popen(
                params,
                cwd=self.work_dir,
                stdin=subprocess.PIPE,
                stdout=output_file,
                stderr=subprocess.PIPE,
                text=True
            )
            self.process = process
            return process

    def create_file(self, filename, content):
        with open(os.path.join(self.work_dir, filename), "w") as f:
            f.write(content)
            
    def flush(self):
        try:
            self.process.stdin.flush()
        except BrokenPipeError as e:
            thread_print(f"Node {self.identifier} : {e}")

    def quit(self):
        self.process.stdin.write("QUIT\n")
        self.flush()
        self.process.wait()

    def print_state(self):
        self.process.stdin.write("PRINTSTATE\n")
        self.flush()

    def store_file(self, file_path):
        self.process.stdin.write("STOREFILE\n")
        self.flush()
        time.sleep(1)
        self.process.stdin.write(file_path + "\n")
        self.flush()
        thread_print(f"{self.process_name} {self.bind_port} stored file {file_path}")

    def get_file(self, file_path):
        self.process.stdin.write("GETFILE\n")
        self.flush()
        time.sleep(1)
        self.process.stdin.write(file_path + "\n")
        self.flush()
        thread_print(f"{self.process_name} {self.bind_port} got file {file_path}")

    def check_file_download(self, file_path, file_content):
        try:
            with open(os.path.join(self.work_dir + '/download/', file_path), "r") as f:
                content = f.read()
                if content == file_content:
                    thread_print(f"{self.process_name} {self.bind_port} downloaded file {file_path} successfully")
                else:
                    thread_print(f"{self.process_name} {self.bind_port} downloaded file {file_path} unsuccessfully")
                    os._exit(1)
        except FileNotFoundError:
            thread_print(f"{self.process_name} {self.bind_port} downloaded file {file_path} unsuccessfully")
            os._exit(1)
