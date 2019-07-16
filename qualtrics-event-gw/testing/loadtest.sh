#!/bin/bash


echo "What is the target load (rps)?"
read rps
echo "What are the stage increments?"
read stage
echo "What is the warmup period per stage(minutes)?"
read stage_duration
echo "What is the desired concurrency?"
read concurrency
echo "What should be the request?"
read QUALTRICS_PAYLOAD
echo "What is the Gateway URL?"
read QUALTRICS_GW_URL


warmup_rps=$stage
while [ "$warmup_rps" -lt "$rps" ]
do
    echo "Starting warmup with ${warmup_rps} for ${stage_duration} minute(s)"
    (( times = warmup_rps * stage_duration * 60 ))
    loadtest -c "$concurrency" --rps "$warmup_rps" -n "$times" \
            "$QUALTRICS_GW_URL" \
            -T 'application/x-www-form-urlencoded' -P ''"$QUALTRICS_PAYLOAD"''
    (( warmup_rps = warmup_rps + stage ))
    echo "Ending warmup with $(warmup_rps)"
done

echo "Starting loadtest with ${rps} rps and concurrency ${concurrency}"
loadtest -c "$concurrency" --rps "$rps" \
            "$QUALTRICS_GW_URL" \
            -T 'application/x-www-form-urlencoded' -P ''"$QUALTRICS_PAYLOAD"''




