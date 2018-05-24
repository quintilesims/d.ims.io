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
DST_REPO_URI   = os.environ['DST_REPO_URI']

class Repository:
    def __init__(self, name, uri, images):
        self.name = name
        self.uri = uri
        self.images = images

def create_dst_tag(repoName, imageTag):
    return '{0}/{1}:{2}'.format(DST_REPO_URI, repoName, imageTag)

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

def calculate_catalog_diff(src, dst):
    diff = {}
 
    for name, repo in src.iteritems():
        if not name in dst:
            diff[name] = repo
            continue

        image_diff = copy.deepcopy(repo.images)
        for src_image in repo.images:
            src_tag = src_image['imageTag']
            for dst_image in dst[name].images:
                dst_tag = dst_image['imageTag']

                if src_tag == dst_tag:
                    image_diff.remove(src_image)
                    break
                  
        if len(image_diff) > 0:
            diff[name] = Repository(repo.name, repo.uri, image_diff)

    return diff

def migrate(repos, src, dst):
    for repo in repos:
        print('Migrating {0} images(s) for repo {1}'.format(len(repo.images), repo.name))
        
        try:
            dst.create_repository(repositoryName=repo.name)
        except:
            pass

        for image in repo.images:
            src_tag = '{0}:{1}'.format(repo.uri, image['imageTag'])
            dst_tag = create_dst_tag(repo.name, image['imageTag'])
            
            subprocess.call(['docker', 'pull', src_tag])
            subprocess.call(['docker', 'tag', src_tag, dst_tag])
            subprocess.call(['docker', 'push', dst_tag])

def split_work(repos, size):
    return [repos[i::size] for i in xrange(size)]

def main():
    src_client = boto3.client('ecr',
        region_name=SRC_REGION,
        aws_access_key_id=SRC_ACCESS_KEY,
        aws_secret_access_key=SRC_SECRET_KEY)

    dst_client = boto3.client('ecr',
        region_name=DST_REGION,
        aws_access_key_id=DST_ACCESS_KEY,
        aws_secret_access_key=DST_SECRET_KEY)
 
    print('Fetching source catalog')
    src_catalog = get_catalog(src_client)

    print('Fetching destination catalog')
    dst_catalog = get_catalog(dst_client)

    # Migrate repos will least amount of images first
    print('Calculating difference')
    diff = calculate_catalog_diff(src_catalog, dst_catalog).values()
    diff = sorted(diff, key=lambda r: len(r.images))

    print('Migrating the following repositories:')
    for repo in diff:
        print('{0} ({1} images)'.format(repo.name, len(repo.images)))

    threads = []
    for repos in split_work(diff, min(multiprocessing.cpu_count(), MAX_WORKERS)):
        thread = threading.Thread(target=migrate, args=(repos, src_client, dst_client))
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
