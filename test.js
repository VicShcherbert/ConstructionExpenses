const walkDog = () => {
    return new Promise((resolve, reject) => {
        const dogWalked = true;
        setTimeout(() => {
            if (dogWalked)
                resolve('Dog walked');
            else
                reject('Failed to walk the dog');
        }, 1000);
    });
}

walkDog();