# ccpalert

CCPAlert is the alerting component of CCPMetrics. It provides a simple threshold based alerting service and can send alerts via Email and PagerDuty. 

### CCPAlertQL
Alerting rules are created via CCPAlertQL, a simple SQL inspired domain specific language. Queries to create alerts take the following form:

```
ALERT <alert name> IF <metric name> <operator> <threshold value> TEXT <description of alert> 
```

The alert name is simply an identifier for the alert. The metric name is the metric which the alert corresponds to. The operator and threshold specify when the alert is triggered. Here are several more concrete examples:

```
ALERT cpuOnFireAlert IF superImportantServer.cpuUsage > 100 TEXT "Critical production server is heavily loaded"
ALERT noplayers IF tq.currentPlayers == 0 TEXT "something has gone badly wrong"
```
