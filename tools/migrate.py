import os
import boto3
import time
import copy
import subprocess
import threading
import multiprocessing

MAX_WORKERS    = 1
SRC_ACCESS_KEY = os.environ['SRC_ACCESS_KEY']
SRC_SECRET_KEY = os.environ['SRC_SECRET_KEY']
SRC_REGION     = os.environ['SRC_REGION']
DST_ACCESS_KEY = os.environ['DST_ACCESS_KEY']
DST_SECRET_KEY = os.environ['DST_SECRET_KEY']
DST_REGION     = os.environ['DST_REGION']

class Repository:
    def __init__(self, name, uri, images):
        self.name = name
        self.uri = uri
        self.images = images

def get_catalog(client):
    catalog = {}

    repos = client.describe_repositories()
    for repo in repos['repositories']:
        name = repo['repositoryName']
        uri = repo['repositoryUri']

        images = []
        paginator = client.get_paginator('list_images')

        for output in paginator.paginate(repositoryName=name, filter={'tagStatus': 'TAGGED'}):
            images = images + output['imageIds']

        catalog[name] = Repository(name, uri, images)

    return catalog

def calculate_catalog_diff(source, dest):
    diff = {}
 
    for name, repo in source.iteritems():
        if not name in dest:
            diff[name] = repo
            continue

        image_diff = copy.deepcopy(repo.images)
        for source_image in repo.images:
            for dest_image in dest[name].images:
                if source_image['imageTag'] == dest_image['imageTag']:
                    image_diff.remove(source_image)
                    break
                  
        if len(image_diff) > 0:
            diff[name] = Repository(repo.name, repo.uri, image_diff)

    return diff

def migrate(repos, source, dest):
    for repo in repos:
        print('Migrating {0} images(s) for repo {1}'.format(len(repo.images), repo.name))
        
        try:
            dest.create_repository(repositoryName=repo.name)
        except:
            pass

        for image in repo.images:
            source_tag = '{0}:{1}'.format(repo.uri, image['imageTag'])
            dest_tag = source_tag.replace(source._client_config.region_name, dest._client_config.region_name)
  
            subprocess.call(['docker', 'pull', source_tag])
            subprocess.call(['docker', 'tag', source_tag, dest_tag])
            subprocess.call(['docker', 'push', dest_tag])

        subprocess.call(['bash', '-c', 'docker rmi -f $(docker images -q)'])

def split_work(repos, size):
    return [repos[i::size] for i in xrange(size)]

def main():
    source_client = boto3.client('ecr',
        region_name=SRC_REGION,
        aws_access_key_id=SRC_ACCESS_KEY,
        aws_secret_access_key=SRC_SECRET_KEY)

    dest_client = boto3.client('ecr',
        region_name=DST_REGION,
        aws_access_key_id=DST_ACCESS_KEY,
        aws_secret_access_key=DST_SECRET_KEY)
 
    print('Fetching source catalog')
    source_catalog = get_catalog(source_client)

    print('Fetching destination catalog')
    dest_catalog = get_catalog(dest_client)

    print('Calculating difference')
    diff = calculate_catalog_diff(source_catalog, dest_catalog)

    threads = []
    for repos in split_work(diff.values(), min(multiprocessing.cpu_count(), MAX_WORKERS)):
        thread = threading.Thread(target=migrate, args=(repos, source_client, dest_client))
        thread.daemon = True
        threads.append(thread)

    print('Migrating {0} repos between {1} thread(s)'.format(len(diff), len(threads)))
    time.sleep(2)

    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()

if __name__ == '__main__':
    main()
