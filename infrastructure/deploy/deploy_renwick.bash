


cp -a $CODEBUILD_SRC_DIR_goartifacts/. ~/goarts/

bucket=$(aws cloudformation describe-stacks --stack-name $CODEPIPELINE_STACK_NAME --query 'Stacks[0].Outputs[?OutputKey==`CodeBucket`].OutputValue' --output text)

echo "Bucket: $bucket"

aws cloudformation package \
    --template-file infrastructure/bot/renwick.yaml \
    --s3-bucket $bucket
    --output yaml > renwick_gen.yaml

aws cloudformation deploy \
    --template-file renwick_gen.yaml \
    --stack-name renwick
    --capabilities CAPABILITY_IAM