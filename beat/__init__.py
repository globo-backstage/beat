import pkg_resources

import cliquet
from pyramid.config import Configurator

__version__ = pkg_resources.get_distribution(__package__).version


def main(global_config, **settings):
    """ This function returns a Pyramid WSGI application.
    """
    config = Configurator(settings=settings)
    config.add_settings({'cliquet.project_name': 'beat'})
    cliquet.initialize(config, __version__)
    config.include('pyramid_chameleon')
    # config.add_static_view('static', 'static', cache_max_age=3600)
    # config.add_route('home', '/')
    config.scan("beat.views")

    return config.make_wsgi_app()
