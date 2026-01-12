from ray.train.torch import TorchTrainer
from ray.train import ScalingConfig
import ray
import time
from torch.utils.tensorboard import SummaryWriter
import os
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import TensorDataset, DataLoader
from sklearn.datasets import load_breast_cancer
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler

def train_loop_per_worker(config):
    import ray.train.torch
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    worker_id = ray.worker.global_worker.node.unique_id
    log_dir = os.path.join(config["log_dir"], f"worker_{worker_id}")
    writer = SummaryWriter(log_dir=log_dir)

    data = load_breast_cancer()
    X = data.data
    y = data.target

    scaler = StandardScaler()
    X = scaler.fit_transform(X)

    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)
    X_train = torch.tensor(X_train, dtype=torch.float32).to(device)
    y_train = torch.tensor(y_train, dtype=torch.long).to(device)

    dataset = TensorDataset(X_train, y_train)
    dataloader = DataLoader(dataset, batch_size=32, shuffle=True, num_workers=0)

    model = nn.Sequential(
        nn.Linear(X_train.shape[1], 1000),
        nn.ReLU(),
        nn.Linear(1000, 2)

    ).to(device)

    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=1e-3)

    start_time = time.time()
    for epoch in range(200):
        total_loss = 0
        for batch_X, batch_y in dataloader:
            output = model(batch_X)
            loss = criterion(output, batch_y)
            optimizer.zero_grad()
            loss.backward()
            optimizer.step()
            total_loss += loss.item()

        avg_loss = total_loss / len(dataloader)
        elapsed_time = time.time() - start_time

        print(f"[Worker {worker_id}] Epoch {epoch}, Loss: {avg_loss:.4f}, Time Elapsed: {elapsed_time:.2f}s")
        writer.add_scalar('Loss/train', avg_loss, epoch)
        writer.add_scalar('Time/elapsed_seconds', elapsed_time, epoch)

    writer.close()

if __name__ == "__main__":
    ray.init(ignore_reinit_error=True)

    start = time.time()
    log_dir = "/tmp/ray/tensorboard_logs"

    trainer = TorchTrainer(
        train_loop_per_worker=train_loop_per_worker,
        scaling_config=ScalingConfig(
            num_workers=2,               # 2 workers = 2 GPUs
            use_gpu=True,                # Required for GPU usage
            resources_per_worker={"CPU": 1, "GPU": 1},
            placement_strategy="SPREAD"  # Optional: spread across nodes
        ),
        train_loop_config={"log_dir": log_dir}
    )

    result = trainer.fit()

    print(f"[Ray TorchTrainer - GPU] Done in {time.time() - start:.2f} seconds")

    ray.shutdown()
