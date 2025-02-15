from diagrams import Diagram
from diagrams.aws.compute import EC2
from diagrams.aws.database import Aurora

with Diagram("TIDE/DIP Service", show=False):
    identity = EC2("Identity")
    transaction = EC2("Transaction")
    portfolio = EC2("Portfolio")
    advisor = EC2("Advisor Dashbaord")
    aurora = Aurora("DIP Aurora")
    user_mgmt = EC2("User Profiling")

    identity >> user_mgmt
    transaction >> user_mgmt
    portfolio >> user_mgmt
    advisor >> user_mgmt
    user_mgmt >> aurora
