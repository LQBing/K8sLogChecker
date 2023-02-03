#!/bin/bash  
  
# resource type
while [[ -z "$KLC_RESOURCE_TYPE" ]];  
do
echo "[required] Enter the resource type(pod/deployment/statfulset):"  
read KLC_RESOURCE_TYPE
if [[ "pod" != "$KLC_RESOURCE_TYPE" ]] && [[ "deployment" != "$KLC_RESOURCE_TYPE" ]] && [[ "statfulset" != "$KLC_RESOURCE_TYPE" ]];then
KLC_RESOURCE_TYPE=""
fi
done
# namespace
echo "Enter the namespace(default value is 'default'):"
read KLC_NAMESPACE
if [[ -z "$KLC_NAMESPACE"  ]]; then
KLC_NAMESPACE="default"
fi
# resource name
while [[ -z "$KLC_RESOURCE_NAME" ]];  
do
echo "[required] Enter the resource name:"
read KLC_RESOURCE_NAME
done
# scan log count
echo "Enter the log lines you want scan(default value is 20):"
read KLC_SCAN_LOG_COUNT
if [[ -z "$KLC_SCAN_LOG_COUNT"  ]]; then
KLC_SCAN_LOG_COUNT="20"
fi
# error content
while [[ -z "$KLC_ERR_CONTENT"  ]];  
do
echo "[required] Enter the error strings you want alert in logs(for example 'ElasticsearchError|[error]'):"  
read KLC_ERR_CONTENT
done
# ignore content
echo "Enter the ignore strings you want alert in logs(for example 'deprecated'):"  
read KLC_IGNORE_CONTENT
# send email
while [[ -z "$KLC_SEND_MAIL" ]];  
do
echo "[required] Do you want send email when alert(yes/no):"
read KLC_SEND_MAIL
if [[ "yes" != "$KLC_SEND_MAIL" ]] && [[ "no" != "$KLC_SEND_MAIL" ]];then
KLC_SEND_MAIL=""
fi
done
# input mail vars if send mail
if [[ "yes" == "$KLC_SEND_MAIL" ]];then
    # mail from
    while [[ -z "$MAIL_FROM" ]];  
    do
    echo "[required] Enter the email from(for example 'abc@mail.com'):"
    read MAIL_FROM
    done
    # mail password
    while [[ -z "$MAIL_PASSWORD" ]];  
    do
    echo "[required] Enter the email password:"
    read MAIL_PASSWORD
    done
    # mail host
    while [[ -z "$MAIL_HOST" ]];  
    do
    echo "[required] Enter the email host(for example 'smtp.office365.com'):"
    read MAIL_HOST
    done
    # mail port
    while [[ -z "$MAIL_PORT" ]];  
    do
    echo "[required] Enter the email smtp server port (for example 587):"
    read MAIL_PORT
    done
    # mail to
    while [[ -z "$MAIL_TO" ]];  
    do
    echo "[required] Enter the email send to (for example rec@mail.com;cc@mail.com):"
    read MAIL_TO
    done
    # a string in the end of mail subject
    echo "Enter a string in the end of mail subject:"
    read MAIL_ADD_SUBJECT
fi

cp k8s.yaml.example k8s.yaml
sed -i "s/namespace: default/namespace: \"$KLC_NAMESPACE\"/g" k8s.yaml
sed -i "s/KLC_NAMESPACE: .*/KLC_NAMESPACE: \"$(echo -ne "$KLC_NAMESPACE"|base64)\"/g" k8s.yaml
sed -i "s/KLC_SCAN_LOG_COUNT: .*/KLC_SCAN_LOG_COUNT: \"$(echo -ne "$KLC_SCAN_LOG_COUNT"|base64)\"/g" k8s.yaml
sed -i "s/KLC_RESOURCE_TYPE: .*/KLC_RESOURCE_TYPE: \"$(echo -ne "$KLC_RESOURCE_TYPE"|base64)\"/g" k8s.yaml
sed -i "s/KLC_RESOURCE_NAME: .*/KLC_RESOURCE_NAME: \"$(echo -ne "$KLC_RESOURCE_NAME"|base64)\"/g" k8s.yaml
sed -i "s/KLC_ERR_CONTENT: .*/KLC_ERR_CONTENT: \"$(echo -ne "$KLC_ERR_CONTENT"|base64)\"/g" k8s.yaml
sed -i "s/KLC_IGNORE_CONTENT: .*/KLC_IGNORE_CONTENT: \"$(echo -ne "$KLC_IGNORE_CONTENT"|base64)\"/g" k8s.yaml

if [[ "yes" == "$KLC_SEND_MAIL" ]];then
    sed -i "s/KLC_SEND_MAIL: .*/KLC_SEND_MAIL: \"$(echo -ne "true"|base64)\"/g" k8s.yaml
    sed -i "s/MAIL_FROM: .*/MAIL_FROM: \"$(echo -ne "$MAIL_FROM"|base64)\"/g" k8s.yaml
    sed -i "s/MAIL_PASSWORD: .*/MAIL_PASSWORD: \"$(echo -ne "$MAIL_PASSWORD"|base64)\"/g" k8s.yaml
    sed -i "s/MAIL_HOST: .*/MAIL_HOST: \"$(echo -ne "$MAIL_HOST"|base64)\"/g" k8s.yaml
    sed -i "s/MAIL_TO: .*/MAIL_TO: \"$(echo -ne "$MAIL_TO"|base64)\"/g" k8s.yaml
    sed -i "s/MAIL_ADD_SUBJECT: .*/MAIL_ADD_SUBJECT: \"$(echo -ne "$MAIL_ADD_SUBJECT"|base64)\"/g" k8s.yaml
else
    sed -i "s/KLC_SEND_MAIL: .*//g" k8s.yaml
    sed -i "s/MAIL_FROM: .*//g" k8s.yaml
    sed -i "s/MAIL_PASSWORD: .*//g" k8s.yaml
    sed -i "s/MAIL_HOST: .*//g" k8s.yaml
    sed -i "s/MAIL_TO: .*//g" k8s.yaml
    sed -i "s/MAIL_ADD_SUBJECT: .*//g" k8s.yaml
fi


echo "====="
declare -p |grep KLC_
declare -p |grep MAIL_
