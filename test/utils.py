import hashlib
import threading

mod = 64


def generate_hash(elt: str) -> int:
    sha1 = hashlib.sha1()
    sha1.update(elt.encode('utf-8'))
    hash_bytes = sha1.digest()
    return int.from_bytes(hash_bytes, byteorder='big') % mod


print_lock = threading.Lock()


def thread_print(*args, **kwargs):
    with print_lock:
        thread_name = threading.current_thread().name
        print(f"[{thread_name}] ", *args, **kwargs)
