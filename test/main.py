import threading
import time
from random import Random

from model import Node
from utils import mod
from utils import thread_print

random = Random()

def create_nodes(node_num, start_port, filename):
    nodes = []
    for i in range(node_num):
        node = Node(
            f"./node_" + str(start_port + i),
            f"node{i}", "127.0.0.1", str(start_port + i)
        )
        node.create_file(filename, f"This is file.txt from node{i}")
        nodes.append(node)
    return nodes

def stabilize_network_normal(nodes):
    nodes[0].create()
    time.sleep(10)
    for i in range(1, len(nodes)):
        nodes[i].join(nodes[0])
    time.sleep(50)  # enough time to stabilize the network
    
def stabilize_network_join(nodes):
    nodes[0].create()
    time.sleep(10)
    for i in range(1, len(nodes)):
        time.sleep(10)
        nodes[i].join(nodes[random.randint(0, i - 1)])
    time.sleep(50)  # enough time to stabilize the network
    
def print_all(nodes):
    for i in range(0, len(nodes)):
        nodes[i].print_state()

def perform_file_operations(nodes, filename):
    store_node = random.randint(0, len(nodes) - 1)
    get_node = random.randint(0, len(nodes) - 1)
    thread_print(f"Store node: {store_node}, Get node: {get_node}")
    nodes[store_node].store_file(filename)
    time.sleep(10)
    nodes[get_node].get_file(filename)
    time.sleep(10)
    nodes[get_node].check_file_download(filename, "This is file.txt from node" + str(store_node))

def quit_nodes(nodes):
    for i in range(1, len(nodes)):
        nodes[i].quit()
    time.sleep(5)
    nodes[0].quit()

def normal_test(node_num, start_port):
    filename = "file_" + str(random.random() % mod) + ".txt"
    nodes = create_nodes(node_num, start_port, filename)
    stabilize_network_normal(nodes)
    print_all(nodes)
    perform_file_operations(nodes, filename)
    quit_nodes(nodes)
    thread_print("Normal test passed")

def join_test(node_num, start_port):
    filename = "file_" + str(random.random() % mod) + ".txt"
    nodes = create_nodes(node_num, start_port, filename)
    stabilize_network_join(nodes)
    print_all(nodes)
    perform_file_operations(nodes, filename)
    quit_nodes(nodes)
    thread_print("Join test passed")

if __name__ == "__main__":
    threads = []
    start_port = 4170
    node_num = 10
    for i in range(3):
        join_thread = threading.Thread(target=join_test, args=(node_num, start_port), name="Thread join_" + str(i))
        join_thread.start()
        threads.append(join_thread)
        start_port += node_num + 1
        normal_thread = threading.Thread(target=normal_test, args=(node_num, start_port), name="Thread normal_" + str(i))
        normal_thread.start()
        threads.append(normal_thread)
        start_port += node_num + 1
    for join_thread in threads:
        join_thread.join()
    print("All tests passed")