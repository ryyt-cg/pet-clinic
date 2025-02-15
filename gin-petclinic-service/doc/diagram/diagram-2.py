from diagrams import Diagram
from diagrams.aws.compute import EC2
from diagrams.aws.database import Aurora

with Diagram("TIDE/DIP Service 2", show=False):
    identity = EC2("Identity")
    transaction = EC2("Transaction")
    portfolio = EC2("Portfolio")
    advisor = EC2("Advisor Dashbaord")
    aurora = Aurora("DIP Aurora")

    identity >> aurora
    transaction >> aurora
    portfolio >> aurora
    advisor >> aurora