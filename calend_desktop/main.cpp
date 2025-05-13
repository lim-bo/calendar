#include "mainwindow.h"

#include <QApplication>
#include "cfg.h"
int main(int argc, char *argv[])
{
    load();
    QApplication a(argc, argv);
    MainWindow w;
    w.show();
    return a.exec();
}
