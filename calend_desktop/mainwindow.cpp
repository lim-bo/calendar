#include "mainwindow.h"
#include "ui_mainwindow.h"
MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow), event(nullptr)

{
    ui->setupUi(this);
    bool fl = false;
    QString uidProvided;
    auth auth_window(&fl, &uidProvided);
    auth_window.exec();
    if (!fl) {
        exit(0);
    }
    this->uid = uidProvided;

    ui->frame->setLayout(new QVBoxLayout());
}

MainWindow::~MainWindow()
{
    delete ui;
    delete event;
}

void MainWindow::on_prof_button_2_clicked()
{
    profile prof(uid);
    prof.exec();

}


void MainWindow::on_mycal_clicked()
{
    int type = ui->caltype->currentIndex();

    QLayoutItem* item;
    while ((item = ui->frame->layout()->takeAt(0)) != nullptr) {
        delete item->widget();
        delete item;
    }
    if (type == 0){
        Cal_30* cal30 = new Cal_30(uid, this);


        ui->frame->layout()->addWidget(cal30);


        cal30->show();

    }
    if (type == 1){
        cal_7* cal7 = new cal_7(uid, this);
        ui->frame->layout()->addWidget(cal7);
        cal7->show();


    }
    if (type == 2){
        cal_1* cal1 = new cal_1(uid, this);
        ui->frame->layout()->addWidget(cal1);
        cal1->show();

    }


}

void MainWindow::on_Event_push_clicked()
{
    Event *ev = new Event(uid);
    ev->setAttribute(Qt::WA_DeleteOnClose);
    ev->show();
}

void MainWindow::on_my_events_clicked()
{
    QLayoutItem* item;
    while ((item = ui->frame->layout()->takeAt(0)) != nullptr) {
        delete item->widget();
        delete item;
    }
    eventsforme *eventsWindow = new eventsforme(uid, this);
    eventsWindow->setAttribute(Qt::WA_DeleteOnClose);
    ui->frame->layout()->addWidget(eventsWindow);
    eventsWindow->show();

}

