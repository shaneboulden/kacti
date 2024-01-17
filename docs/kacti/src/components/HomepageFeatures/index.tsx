import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Kubernetes native',
    description: (
      <>
        Built from the ground up to support Kubernetes and Kubernetes-native
        security platforms, like StackRox.
      </>
    ),
  },
  {
    title: 'Human readable',
    description: (
      <>
        Simple, human-readable testing formats, allowing you to easily describe, customize and 
        share tests.
      </>
    ),
  },
  {
    title: 'Integrate with CI/CD',
    description: (
      <>
        Easily integrate with CI/CD pipelines to functionally verify 
        Kubernetes security policy changes
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
