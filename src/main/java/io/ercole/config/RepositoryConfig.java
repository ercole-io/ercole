package io.ercole.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.data.rest.core.config.RepositoryRestConfiguration;
import org.springframework.data.rest.webmvc.config.RepositoryRestConfigurerAdapter;

import io.ercole.model.Alert;
import io.ercole.model.CurrentHost;
import io.ercole.model.License;

/**
 * RepositoryConfig custom configuration.
 */
@Configuration
public class RepositoryConfig extends RepositoryRestConfigurerAdapter {

    /**
     * Public method.
     *
     * @param config the RepositoryRestConfiguration.
     */
    @Override
    public void configureRepositoryRestConfiguration(final RepositoryRestConfiguration config) {
        config.exposeIdsFor(
                CurrentHost.class,
                Alert.class,
                License.class
        );
    }
}
